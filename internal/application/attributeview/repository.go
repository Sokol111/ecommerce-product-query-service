package attributeview

import "context"

// Repository defines the interface for attribute view persistence
type Repository interface {
	// Upsert inserts or updates an attribute if the incoming version is greater than stored
	Upsert(ctx context.Context, attribute *AttributeView) error

	// FindByIDs returns attributes by their IDs
	FindByIDs(ctx context.Context, ids []string) ([]*AttributeView, error)
}
