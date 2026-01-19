package client

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	catalogapi "github.com/Sokol111/ecommerce-catalog-service-api/gen/httpapi"
)

// AttributeClient provides access to catalog-service for attributes
type AttributeClient interface {
	// GetAttributeByID fetches attribute data by ID
	GetAttributeByID(ctx context.Context, id string) (*catalogapi.AttributeResponse, error)
	// GetAttributesByIDs fetches multiple attributes by their IDs
	GetAttributesByIDs(ctx context.Context, ids []string) (map[string]*catalogapi.AttributeResponse, error)
}

type attributeClient struct {
	client *catalogapi.Client
}

// newAttributeClient creates a new attribute client
func newAttributeClient(client *catalogapi.Client) AttributeClient {
	return &attributeClient{client: client}
}

// GetAttributeByID fetches a single attribute by ID
func (c *attributeClient) GetAttributeByID(ctx context.Context, id string) (*catalogapi.AttributeResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid attribute ID %s: %w", id, err)
	}

	res, err := c.client.GetAttributeById(ctx, catalogapi.GetAttributeByIdParams{ID: uid})
	if err != nil {
		return nil, fmt.Errorf("failed to fetch attribute %s: %w", id, err)
	}

	switch v := res.(type) {
	case *catalogapi.AttributeResponse:
		return v, nil
	case *catalogapi.GetAttributeByIdNotFound:
		return nil, fmt.Errorf("attribute %s not found", id)
	default:
		return nil, fmt.Errorf("unexpected response type for attribute %s", id)
	}
}

// GetAttributesByIDs fetches multiple attributes by their IDs
func (c *attributeClient) GetAttributesByIDs(ctx context.Context, ids []string) (map[string]*catalogapi.AttributeResponse, error) {
	result := make(map[string]*catalogapi.AttributeResponse, len(ids))

	// TODO: Consider adding batch endpoint to catalog-service for better performance
	for _, id := range ids {
		attr, err := c.GetAttributeByID(ctx, id)
		if err != nil {
			return nil, err
		}
		result[id] = attr
	}

	return result, nil
}
