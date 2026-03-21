package mongoaudit

import (
	"context"
	"encoding/json"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	financedomain "vct-platform/backend/internal/modules/finance/domain"
)

// Repository stores inspection-grade void snapshots in MongoDB.
type Repository struct {
	collection *mongo.Collection
}

// NewRepository constructs the MongoDB audit adapter.
func NewRepository(collection *mongo.Collection) *Repository {
	return &Repository{collection: collection}
}

// RecordVoid persists the before/after snapshots around a void operation.
func (r *Repository) RecordVoid(ctx context.Context, audit financedomain.VoidAuditLog) error {
	if r == nil || r.collection == nil {
		return fmt.Errorf("mongo audit collection is not configured")
	}

	before, err := normalizeJSONDocument(audit.Before)
	if err != nil {
		return fmt.Errorf("normalize void audit before snapshot: %w", err)
	}
	after, err := normalizeJSONDocument(audit.After)
	if err != nil {
		return fmt.Errorf("normalize void audit after snapshot: %w", err)
	}

	_, err = r.collection.InsertOne(ctx, bson.M{
		"_id":               audit.ID,
		"company_code":      audit.CompanyCode,
		"original_entry_id": audit.OriginalEntryID,
		"reversal_entry_id": audit.ReversalEntryID,
		"actor_id":          audit.ActorID,
		"reason":            audit.Reason,
		"voided_at":         audit.VoidedAt,
		"before":            before,
		"after":             after,
	})
	if err != nil {
		return fmt.Errorf("insert void audit log %s: %w", audit.ID, err)
	}

	return nil
}

func normalizeJSONDocument(value any) (any, error) {
	if value == nil {
		return bson.M{}, nil
	}

	raw, err := json.Marshal(value)
	if err != nil {
		return nil, err
	}

	var normalized any
	if err := json.Unmarshal(raw, &normalized); err != nil {
		return nil, err
	}
	return normalized, nil
}
