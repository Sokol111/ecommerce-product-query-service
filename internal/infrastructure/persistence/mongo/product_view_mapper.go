package mongo

import (
	commonsmongo "github.com/Sokol111/ecommerce-commons/pkg/persistence/mongo"
	"github.com/Sokol111/ecommerce-product-query-service/internal/domain/productview"
)

type productViewMapper struct{}

func newProductViewMapper() *productViewMapper {
	return &productViewMapper{}
}

func (m *productViewMapper) ToEntity(domain *productview.ProductView) *productViewEntity {
	return &productViewEntity{
		ID:          domain.ID,
		Version:     domain.Version,
		Name:        domain.Name,
		Description: domain.Description,
		Price:       domain.Price,
		Quantity:    domain.Quantity,
		ImageID:     domain.ImageID,
		Enabled:     domain.Enabled,
		CreatedAt:   domain.CreatedAt,
		ModifiedAt:  domain.ModifiedAt,
	}
}

func (m *productViewMapper) ToDomain(entity *productViewEntity) *productview.ProductView {
	return productview.Reconstruct(
		entity.ID,
		entity.Version,
		entity.Name,
		entity.Description,
		entity.Price,
		entity.Quantity,
		entity.ImageID,
		entity.Enabled,
		entity.CreatedAt,
		entity.ModifiedAt,
	)
}

func (m *productViewMapper) GetID(entity *productViewEntity) string {
	return entity.ID
}

func (m *productViewMapper) GetVersion(entity *productViewEntity) int {
	return entity.Version
}

func (m *productViewMapper) SetVersion(entity *productViewEntity, version int) {
	entity.Version = version
}

// Ensure mapper implements the interface
var _ commonsmongo.EntityMapper[productview.ProductView, productViewEntity] = (*productViewMapper)(nil)
