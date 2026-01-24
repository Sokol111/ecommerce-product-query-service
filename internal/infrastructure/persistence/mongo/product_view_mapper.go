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
	var attributes []productAttributeEntity
	if len(domain.Attributes) > 0 {
		attributes = make([]productAttributeEntity, len(domain.Attributes))
		for i, attr := range domain.Attributes {
			attributes[i] = productAttributeEntity{
				AttributeID:      attr.AttributeID,
				Slug:             attr.Slug,
				OptionSlugValue:  attr.OptionSlugValue,
				OptionSlugValues: attr.OptionSlugValues,
				NumericValue:     attr.NumericValue,
				TextValue:        attr.TextValue,
				BooleanValue:     attr.BooleanValue,
			}
		}
	}

	return &productViewEntity{
		ID:            domain.ID,
		Version:       domain.Version,
		Name:          domain.Name,
		Description:   domain.Description,
		Price:         domain.Price,
		Quantity:      domain.Quantity,
		ImageID:       domain.ImageID,
		SmallImageURL: domain.SmallImageURL,
		LargeImageURL: domain.LargeImageURL,
		CategoryID:    domain.CategoryID,
		Enabled:       domain.Enabled,
		CreatedAt:     domain.CreatedAt,
		ModifiedAt:    domain.ModifiedAt,
		Attributes:    attributes,
		Attrs:         domain.Attrs,
	}
}

func (m *productViewMapper) ToDomain(entity *productViewEntity) *productview.ProductView {
	var attributes []productview.AttributeValue
	if len(entity.Attributes) > 0 {
		attributes = make([]productview.AttributeValue, len(entity.Attributes))
		for i, attr := range entity.Attributes {
			attributes[i] = productview.AttributeValue{
				AttributeID:      attr.AttributeID,
				Slug:             attr.Slug,
				OptionSlugValue:  attr.OptionSlugValue,
				OptionSlugValues: attr.OptionSlugValues,
				NumericValue:     attr.NumericValue,
				TextValue:        attr.TextValue,
				BooleanValue:     attr.BooleanValue,
			}
		}
	}

	return productview.Reconstruct(
		entity.ID,
		entity.Version,
		entity.Name,
		entity.Description,
		entity.Price,
		entity.Quantity,
		entity.ImageID,
		entity.SmallImageURL,
		entity.LargeImageURL,
		entity.CategoryID,
		entity.Enabled,
		entity.CreatedAt,
		entity.ModifiedAt,
		attributes,
		entity.Attrs,
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
