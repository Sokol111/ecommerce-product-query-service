package client

import (
	"context"
	"fmt"

	attributeapi "github.com/Sokol111/ecommerce-attribute-service-api/gen/httpapi"
)

// AttributeClient provides access to attribute-service
type AttributeClient interface {
	// GetAttributeByID fetches attribute data by ID
	GetAttributeByID(ctx context.Context, id string) (*attributeapi.AttributeResponse, error)
	// GetAttributesByIDs fetches multiple attributes by their IDs
	GetAttributesByIDs(ctx context.Context, ids []string) (map[string]*attributeapi.AttributeResponse, error)
}

type attributeClient struct {
	client *attributeapi.Client
}

// newAttributeClient creates a new attribute client
func newAttributeClient(client *attributeapi.Client) AttributeClient {
	return &attributeClient{client: client}
}

// GetAttributeByID fetches a single attribute by ID
func (c *attributeClient) GetAttributeByID(ctx context.Context, id string) (*attributeapi.AttributeResponse, error) {
	res, err := c.client.GetAttributeById(ctx, attributeapi.GetAttributeByIdParams{ID: id})
	if err != nil {
		return nil, fmt.Errorf("failed to fetch attribute %s: %w", id, err)
	}

	switch v := res.(type) {
	case *attributeapi.AttributeResponse:
		return v, nil
	case *attributeapi.GetAttributeByIdNotFound:
		return nil, fmt.Errorf("attribute %s not found", id)
	default:
		return nil, fmt.Errorf("unexpected response type for attribute %s", id)
	}
}

// GetAttributesByIDs fetches multiple attributes by their IDs
func (c *attributeClient) GetAttributesByIDs(ctx context.Context, ids []string) (map[string]*attributeapi.AttributeResponse, error) {
	result := make(map[string]*attributeapi.AttributeResponse, len(ids))

	// TODO: Consider adding batch endpoint to attribute-service for better performance
	for _, id := range ids {
		attr, err := c.GetAttributeByID(ctx, id)
		if err != nil {
			return nil, err
		}
		result[id] = attr
	}

	return result, nil
}
