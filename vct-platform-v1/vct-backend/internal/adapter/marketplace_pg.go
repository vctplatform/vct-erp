package adapter

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"vct-platform/backend/internal/domain/marketplace"
)

type PgMarketplaceProductRepo struct {
	db *sql.DB
}

func NewPgMarketplaceProductRepo(db *sql.DB) *PgMarketplaceProductRepo {
	return &PgMarketplaceProductRepo{db: db}
}

const marketplaceProductCols = `
id, slug, seller_id, seller_name, seller_role, title, short_description, description,
category, condition, martial_art, price_vnd, compare_at_price_vnd, currency,
stock_quantity, minimum_order_quantity, status, location, featured,
images, tags, specs, shipping, created_at, updated_at
`

func (r *PgMarketplaceProductRepo) List(ctx context.Context) ([]marketplace.Product, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT `+marketplaceProductCols+` FROM marketplace_products ORDER BY featured DESC, updated_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]marketplace.Product, 0)
	for rows.Next() {
		item, scanErr := scanMarketplaceProduct(rows)
		if scanErr != nil {
			return nil, scanErr
		}
		items = append(items, *item)
	}
	return items, rows.Err()
}

func (r *PgMarketplaceProductRepo) GetByID(ctx context.Context, id string) (*marketplace.Product, error) {
	item, err := scanMarketplaceProduct(r.db.QueryRowContext(ctx,
		`SELECT `+marketplaceProductCols+` FROM marketplace_products WHERE id=$1`, id))
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("product not found: %s", id)
	}
	return item, err
}

func (r *PgMarketplaceProductRepo) GetBySlug(ctx context.Context, slug string) (*marketplace.Product, error) {
	item, err := scanMarketplaceProduct(r.db.QueryRowContext(ctx,
		`SELECT `+marketplaceProductCols+` FROM marketplace_products WHERE slug=$1`, slug))
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("product not found: %s", slug)
	}
	return item, err
}

func (r *PgMarketplaceProductRepo) Create(ctx context.Context, product marketplace.Product) (*marketplace.Product, error) {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO marketplace_products (
			id, slug, seller_id, seller_name, seller_role, title, short_description, description,
			category, condition, martial_art, price_vnd, compare_at_price_vnd, currency,
			stock_quantity, minimum_order_quantity, status, location, featured,
			images, tags, specs, shipping, created_at, updated_at
		) VALUES (
			$1,$2,$3,$4,$5,$6,$7,$8,
			$9,$10,$11,$12,$13,$14,
			$15,$16,$17,$18,$19,
			$20,$21,$22,$23,$24,$25
		)
	`,
		product.ID, product.Slug, product.SellerID, product.SellerName, product.SellerRole, product.Title,
		product.ShortDescription, product.Description, string(product.Category), string(product.Condition),
		product.MartialArt, product.PriceVND, product.CompareAtPriceVND, product.Currency,
		product.StockQuantity, product.MinimumOrderQuantity, string(product.Status), product.Location, product.Featured,
		mustJSON(product.Images), mustJSON(product.Tags), mustJSON(product.Specs), mustJSON(product.Shipping),
		product.CreatedAt, product.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *PgMarketplaceProductRepo) Update(ctx context.Context, id string, patch map[string]any) (*marketplace.Product, error) {
	current, err := r.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	applyMarketplaceProductPatch(current, patch)

	_, err = r.db.ExecContext(ctx, `
		UPDATE marketplace_products SET
			slug=$2, seller_id=$3, seller_name=$4, seller_role=$5, title=$6,
			short_description=$7, description=$8, category=$9, condition=$10, martial_art=$11,
			price_vnd=$12, compare_at_price_vnd=$13, currency=$14, stock_quantity=$15,
			minimum_order_quantity=$16, status=$17, location=$18, featured=$19,
			images=$20, tags=$21, specs=$22, shipping=$23, updated_at=$24
		WHERE id=$1
	`,
		current.ID, current.Slug, current.SellerID, current.SellerName, current.SellerRole, current.Title,
		current.ShortDescription, current.Description, string(current.Category), string(current.Condition), current.MartialArt,
		current.PriceVND, current.CompareAtPriceVND, current.Currency, current.StockQuantity,
		current.MinimumOrderQuantity, string(current.Status), current.Location, current.Featured,
		mustJSON(current.Images), mustJSON(current.Tags), mustJSON(current.Specs), mustJSON(current.Shipping), current.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return current, nil
}

func (r *PgMarketplaceProductRepo) Delete(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM marketplace_products WHERE id=$1`, id)
	return err
}

func (r *PgMarketplaceProductRepo) ListBySeller(ctx context.Context, sellerID string) ([]marketplace.Product, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT `+marketplaceProductCols+` FROM marketplace_products WHERE seller_id=$1 ORDER BY updated_at DESC`, sellerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]marketplace.Product, 0)
	for rows.Next() {
		item, scanErr := scanMarketplaceProduct(rows)
		if scanErr != nil {
			return nil, scanErr
		}
		items = append(items, *item)
	}
	return items, rows.Err()
}

type PgMarketplaceOrderRepo struct {
	db *sql.DB
}

func NewPgMarketplaceOrderRepo(db *sql.DB) *PgMarketplaceOrderRepo {
	return &PgMarketplaceOrderRepo{db: db}
}

const marketplaceOrderCols = `
id, order_code, seller_id, seller_name, buyer_name, buyer_phone, buyer_email,
buyer_address, notes, status, payment_status, subtotal_vnd, shipping_fee_vnd,
discount_vnd, total_vnd, created_at, updated_at
`

func (r *PgMarketplaceOrderRepo) List(ctx context.Context) ([]marketplace.Order, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT `+marketplaceOrderCols+` FROM marketplace_orders ORDER BY created_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanOrdersWithItems(ctx, rows)
}

func (r *PgMarketplaceOrderRepo) GetByID(ctx context.Context, id string) (*marketplace.Order, error) {
	order, err := r.scanMarketplaceOrder(ctx, r.db.QueryRowContext(ctx,
		`SELECT `+marketplaceOrderCols+` FROM marketplace_orders WHERE id=$1`, id))
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("order not found: %s", id)
	}
	return order, err
}

func (r *PgMarketplaceOrderRepo) Create(ctx context.Context, order marketplace.Order) (*marketplace.Order, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	_, err = tx.ExecContext(ctx, `
		INSERT INTO marketplace_orders (
			id, order_code, seller_id, seller_name, buyer_name, buyer_phone, buyer_email,
			buyer_address, notes, status, payment_status, subtotal_vnd, shipping_fee_vnd,
			discount_vnd, total_vnd, created_at, updated_at
		) VALUES (
			$1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17
		)
	`,
		order.ID, order.OrderCode, order.SellerID, order.SellerName, order.BuyerName, order.BuyerPhone,
		order.BuyerEmail, order.BuyerAddress, order.Notes, string(order.Status), string(order.PaymentStatus),
		order.SubtotalVND, order.ShippingFeeVND, order.DiscountVND, order.TotalVND, order.CreatedAt, order.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	for _, item := range order.Items {
		_, err = tx.ExecContext(ctx, `
			INSERT INTO marketplace_order_items (
				id, order_id, product_id, product_slug, product_title, unit_price_vnd, quantity, line_total_vnd
			) VALUES ($1,$2,$3,$4,$5,$6,$7,$8)
		`,
			item.ID, order.ID, item.ProductID, item.ProductSlug, item.ProductTitle, item.UnitPriceVND, item.Quantity, item.LineTotalVND,
		)
		if err != nil {
			return nil, err
		}
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}
	return &order, nil
}

func (r *PgMarketplaceOrderRepo) Update(ctx context.Context, id string, patch map[string]any) (*marketplace.Order, error) {
	current, err := r.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	applyMarketplaceOrderPatch(current, patch)

	_, err = r.db.ExecContext(ctx, `
		UPDATE marketplace_orders SET
			status=$2, payment_status=$3, notes=$4, updated_at=$5
		WHERE id=$1
	`, current.ID, string(current.Status), string(current.PaymentStatus), current.Notes, current.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return current, nil
}

func (r *PgMarketplaceOrderRepo) ListBySeller(ctx context.Context, sellerID string) ([]marketplace.Order, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT `+marketplaceOrderCols+` FROM marketplace_orders WHERE seller_id=$1 ORDER BY created_at DESC`, sellerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanOrdersWithItems(ctx, rows)
}

func (r *PgMarketplaceOrderRepo) scanOrdersWithItems(ctx context.Context, rows *sql.Rows) ([]marketplace.Order, error) {
	items := make([]marketplace.Order, 0)
	for rows.Next() {
		order, scanErr := scanMarketplaceOrderRow(rows)
		if scanErr != nil {
			return nil, scanErr
		}
		order.Items, scanErr = r.loadOrderItems(ctx, order.ID)
		if scanErr != nil {
			return nil, scanErr
		}
		items = append(items, *order)
	}
	return items, rows.Err()
}

func (r *PgMarketplaceOrderRepo) scanMarketplaceOrder(ctx context.Context, row interface{ Scan(...any) error }) (*marketplace.Order, error) {
	order, err := scanMarketplaceOrderRow(row)
	if err != nil {
		return nil, err
	}
	order.Items, err = r.loadOrderItems(ctx, order.ID)
	if err != nil {
		return nil, err
	}
	return order, nil
}

func (r *PgMarketplaceOrderRepo) loadOrderItems(ctx context.Context, orderID string) ([]marketplace.OrderItem, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, product_id, product_slug, product_title, unit_price_vnd, quantity, line_total_vnd
		FROM marketplace_order_items
		WHERE order_id=$1
		ORDER BY id ASC
	`, orderID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]marketplace.OrderItem, 0)
	for rows.Next() {
		var item marketplace.OrderItem
		if err := rows.Scan(
			&item.ID, &item.ProductID, &item.ProductSlug, &item.ProductTitle,
			&item.UnitPriceVND, &item.Quantity, &item.LineTotalVND,
		); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func scanMarketplaceProduct(row interface{ Scan(...any) error }) (*marketplace.Product, error) {
	var (
		item         marketplace.Product
		imagesJSON   []byte
		tagsJSON     []byte
		specsJSON    []byte
		shippingJSON []byte
	)
	err := row.Scan(
		&item.ID, &item.Slug, &item.SellerID, &item.SellerName, &item.SellerRole, &item.Title,
		&item.ShortDescription, &item.Description, &item.Category, &item.Condition, &item.MartialArt,
		&item.PriceVND, &item.CompareAtPriceVND, &item.Currency, &item.StockQuantity, &item.MinimumOrderQuantity,
		&item.Status, &item.Location, &item.Featured, &imagesJSON, &tagsJSON, &specsJSON, &shippingJSON,
		&item.CreatedAt, &item.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	if err := decodeJSONColumn(imagesJSON, &item.Images); err != nil {
		return nil, err
	}
	if err := decodeJSONColumn(tagsJSON, &item.Tags); err != nil {
		return nil, err
	}
	if err := decodeJSONColumn(specsJSON, &item.Specs); err != nil {
		return nil, err
	}
	if err := decodeJSONColumn(shippingJSON, &item.Shipping); err != nil {
		return nil, err
	}
	return &item, nil
}

func scanMarketplaceOrderRow(row interface{ Scan(...any) error }) (*marketplace.Order, error) {
	var item marketplace.Order
	err := row.Scan(
		&item.ID, &item.OrderCode, &item.SellerID, &item.SellerName, &item.BuyerName, &item.BuyerPhone,
		&item.BuyerEmail, &item.BuyerAddress, &item.Notes, &item.Status, &item.PaymentStatus,
		&item.SubtotalVND, &item.ShippingFeeVND, &item.DiscountVND, &item.TotalVND, &item.CreatedAt, &item.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func decodeJSONColumn[T any](raw []byte, target *T) error {
	if len(raw) == 0 {
		return nil
	}
	return json.Unmarshal(raw, target)
}

func mustJSON(value any) []byte {
	data, err := json.Marshal(value)
	if err != nil {
		return []byte("null")
	}
	return data
}

func applyMarketplaceProductPatch(item *marketplace.Product, patch map[string]any) {
	if value, ok := patch["slug"].(string); ok {
		item.Slug = strings.TrimSpace(value)
	}
	if value, ok := patch["seller_id"].(string); ok {
		item.SellerID = strings.TrimSpace(value)
	}
	if value, ok := patch["seller_name"].(string); ok {
		item.SellerName = strings.TrimSpace(value)
	}
	if value, ok := patch["seller_role"].(string); ok {
		item.SellerRole = strings.TrimSpace(value)
	}
	if value, ok := patch["title"].(string); ok {
		item.Title = strings.TrimSpace(value)
	}
	if value, ok := patch["short_description"].(string); ok {
		item.ShortDescription = strings.TrimSpace(value)
	}
	if value, ok := patch["description"].(string); ok {
		item.Description = strings.TrimSpace(value)
	}
	if value, ok := patch["category"].(string); ok {
		item.Category = marketplace.ProductCategory(strings.TrimSpace(value))
	}
	if value, ok := patch["condition"].(string); ok {
		item.Condition = marketplace.ProductCondition(strings.TrimSpace(value))
	}
	if value, ok := patch["martial_art"].(string); ok {
		item.MartialArt = strings.TrimSpace(value)
	}
	if value, ok := numericToInt64(patch["price_vnd"]); ok {
		item.PriceVND = value
	}
	if value, ok := numericToInt64(patch["compare_at_price_vnd"]); ok {
		item.CompareAtPriceVND = value
	}
	if value, ok := patch["currency"].(string); ok {
		item.Currency = strings.TrimSpace(value)
	}
	if value, ok := numericToInt(patch["stock_quantity"]); ok {
		item.StockQuantity = value
	}
	if value, ok := numericToInt(patch["minimum_order_quantity"]); ok {
		item.MinimumOrderQuantity = value
	}
	if value, ok := patch["status"].(string); ok {
		item.Status = marketplace.ProductStatus(strings.TrimSpace(value))
	}
	if value, ok := patch["location"].(string); ok {
		item.Location = strings.TrimSpace(value)
	}
	if value, ok := patch["featured"].(bool); ok {
		item.Featured = value
	}
	if value, ok := patch["images"].([]string); ok {
		item.Images = value
	}
	if value, ok := patch["tags"].([]string); ok {
		item.Tags = value
	}
	if value, ok := patch["specs"].([]marketplace.ProductSpec); ok {
		item.Specs = value
	}
	if value, ok := patch["shipping"].(marketplace.ShippingProfile); ok {
		item.Shipping = value
	}
	if value, ok := patch["images"]; ok {
		item.Images = decodeStringSlice(value)
	}
	if value, ok := patch["tags"]; ok {
		item.Tags = decodeStringSlice(value)
	}
	if value, ok := patch["specs"]; ok {
		if decoded, err := decodeSpecs(value); err == nil {
			item.Specs = decoded
		}
	}
	if value, ok := patch["shipping"]; ok {
		if decoded, err := decodeShipping(value); err == nil {
			item.Shipping = decoded
		}
	}
	switch typed := patch["updated_at"].(type) {
	case time.Time:
		item.UpdatedAt = typed
	case string:
		if parsed, err := time.Parse(time.RFC3339Nano, strings.TrimSpace(typed)); err == nil {
			item.UpdatedAt = parsed
		}
	}
}

func applyMarketplaceOrderPatch(item *marketplace.Order, patch map[string]any) {
	if value, ok := patch["status"].(string); ok {
		item.Status = marketplace.OrderStatus(strings.TrimSpace(value))
	}
	if value, ok := patch["payment_status"].(string); ok {
		item.PaymentStatus = marketplace.PaymentStatus(strings.TrimSpace(value))
	}
	if value, ok := patch["notes"].(string); ok {
		item.Notes = strings.TrimSpace(value)
	}
	switch typed := patch["updated_at"].(type) {
	case time.Time:
		item.UpdatedAt = typed
	case string:
		if parsed, err := time.Parse(time.RFC3339Nano, strings.TrimSpace(typed)); err == nil {
			item.UpdatedAt = parsed
		}
	}
}

func numericToInt64(value any) (int64, bool) {
	switch typed := value.(type) {
	case int:
		return int64(typed), true
	case int64:
		return typed, true
	case float64:
		return int64(typed), true
	case json.Number:
		parsed, err := typed.Int64()
		return parsed, err == nil
	default:
		return 0, false
	}
}

func numericToInt(value any) (int, bool) {
	switch typed := value.(type) {
	case int:
		return typed, true
	case int64:
		return int(typed), true
	case float64:
		return int(typed), true
	case json.Number:
		parsed, err := typed.Int64()
		return int(parsed), err == nil
	default:
		return 0, false
	}
}

func decodeStringSlice(value any) []string {
	switch typed := value.(type) {
	case []string:
		return typed
	case []any:
		out := make([]string, 0, len(typed))
		for _, item := range typed {
			text, ok := item.(string)
			if ok && strings.TrimSpace(text) != "" {
				out = append(out, strings.TrimSpace(text))
			}
		}
		return out
	default:
		return nil
	}
}

func decodeSpecs(value any) ([]marketplace.ProductSpec, error) {
	switch typed := value.(type) {
	case []marketplace.ProductSpec:
		return typed, nil
	default:
		data, err := json.Marshal(typed)
		if err != nil {
			return nil, err
		}
		var out []marketplace.ProductSpec
		if err := json.Unmarshal(data, &out); err != nil {
			return nil, err
		}
		return out, nil
	}
}

func decodeShipping(value any) (marketplace.ShippingProfile, error) {
	switch typed := value.(type) {
	case marketplace.ShippingProfile:
		return typed, nil
	default:
		data, err := json.Marshal(typed)
		if err != nil {
			return marketplace.ShippingProfile{}, err
		}
		var out marketplace.ShippingProfile
		if err := json.Unmarshal(data, &out); err != nil {
			return marketplace.ShippingProfile{}, err
		}
		return out, nil
	}
}
