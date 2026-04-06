package adapter

import (
	"context"
	"fmt"

	"vct-platform/backend/internal/domain/marketplace"
	"vct-platform/backend/internal/store"
)

type marketplaceProductRepo struct {
	*StoreAdapter[marketplace.Product]
}

func NewMarketplaceProductRepository(ds store.DataStore) marketplace.ProductRepository {
	return &marketplaceProductRepo{
		StoreAdapter: NewStoreAdapter[marketplace.Product](ds, "marketplace_products"),
	}
}

func (r *marketplaceProductRepo) List(ctx context.Context) ([]marketplace.Product, error) {
	return r.StoreAdapter.List()
}

func (r *marketplaceProductRepo) GetByID(ctx context.Context, id string) (*marketplace.Product, error) {
	return r.StoreAdapter.GetByID(id)
}

func (r *marketplaceProductRepo) GetBySlug(ctx context.Context, slug string) (*marketplace.Product, error) {
	items, err := r.List(ctx)
	if err != nil {
		return nil, err
	}
	for _, item := range items {
		if item.Slug == slug {
			product := item
			return &product, nil
		}
	}
	return nil, fmt.Errorf("marketplace product not found: %s", slug)
}

func (r *marketplaceProductRepo) Create(ctx context.Context, product marketplace.Product) (*marketplace.Product, error) {
	return r.StoreAdapter.Create(product)
}

func (r *marketplaceProductRepo) Update(ctx context.Context, id string, patch map[string]any) (*marketplace.Product, error) {
	return r.StoreAdapter.Update(id, patch)
}

func (r *marketplaceProductRepo) Delete(ctx context.Context, id string) error {
	return r.StoreAdapter.Delete(id)
}

func (r *marketplaceProductRepo) ListBySeller(ctx context.Context, sellerID string) ([]marketplace.Product, error) {
	items, err := r.List(ctx)
	if err != nil {
		return nil, err
	}
	filtered := make([]marketplace.Product, 0)
	for _, item := range items {
		if item.SellerID == sellerID {
			filtered = append(filtered, item)
		}
	}
	return filtered, nil
}

type marketplaceOrderRepo struct {
	*StoreAdapter[marketplace.Order]
}

func NewMarketplaceOrderRepository(ds store.DataStore) marketplace.OrderRepository {
	return &marketplaceOrderRepo{
		StoreAdapter: NewStoreAdapter[marketplace.Order](ds, "marketplace_orders"),
	}
}

func (r *marketplaceOrderRepo) List(ctx context.Context) ([]marketplace.Order, error) {
	return r.StoreAdapter.List()
}

func (r *marketplaceOrderRepo) GetByID(ctx context.Context, id string) (*marketplace.Order, error) {
	return r.StoreAdapter.GetByID(id)
}

func (r *marketplaceOrderRepo) Create(ctx context.Context, order marketplace.Order) (*marketplace.Order, error) {
	return r.StoreAdapter.Create(order)
}

func (r *marketplaceOrderRepo) Update(ctx context.Context, id string, patch map[string]any) (*marketplace.Order, error) {
	return r.StoreAdapter.Update(id, patch)
}

func (r *marketplaceOrderRepo) ListBySeller(ctx context.Context, sellerID string) ([]marketplace.Order, error) {
	items, err := r.List(ctx)
	if err != nil {
		return nil, err
	}
	filtered := make([]marketplace.Order, 0)
	for _, item := range items {
		if item.SellerID == sellerID {
			filtered = append(filtered, item)
		}
	}
	return filtered, nil
}
