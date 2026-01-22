package mongo

import (
	"github.com/samber/lo"

	commonsmongo "github.com/Sokol111/ecommerce-commons/pkg/persistence/mongo"
	"github.com/Sokol111/ecommerce-product-query-service/internal/domain/attributeview"
)

type attributeViewMapper struct{}

func newAttributeViewMapper() *attributeViewMapper {
	return &attributeViewMapper{}
}

func (m *attributeViewMapper) ToEntity(domain *attributeview.AttributeView) *attributeViewEntity {
	return &attributeViewEntity{
		ID:         domain.ID,
		Version:    domain.Version,
		Slug:       domain.Slug,
		Name:       domain.Name,
		Type:       string(domain.Type),
		Unit:       domain.Unit,
		Enabled:    domain.Enabled,
		ModifiedAt: domain.ModifiedAt,
		Options: lo.Map(domain.Options, func(opt attributeview.AttributeOption, _ int) attributeOptionEntity {
			return attributeOptionEntity{
				Slug:      opt.Slug,
				Name:      opt.Name,
				ColorCode: opt.ColorCode,
				SortOrder: opt.SortOrder,
			}
		}),
	}
}

func (m *attributeViewMapper) ToDomain(entity *attributeViewEntity) *attributeview.AttributeView {
	return attributeview.Reconstruct(
		entity.ID,
		entity.Version,
		entity.Slug,
		entity.Name,
		attributeview.AttributeType(entity.Type),
		entity.Unit,
		entity.Enabled,
		entity.ModifiedAt,
		lo.Map(entity.Options, func(opt attributeOptionEntity, _ int) attributeview.AttributeOption {
			return attributeview.AttributeOption{
				Slug:      opt.Slug,
				Name:      opt.Name,
				ColorCode: opt.ColorCode,
				SortOrder: opt.SortOrder,
			}
		}),
	)
}

func (m *attributeViewMapper) GetID(entity *attributeViewEntity) string {
	return entity.ID
}

func (m *attributeViewMapper) GetVersion(entity *attributeViewEntity) int {
	return entity.Version
}

func (m *attributeViewMapper) SetVersion(entity *attributeViewEntity, version int) {
	entity.Version = version
}

// Ensure mapper implements the interface
var _ commonsmongo.EntityMapper[attributeview.AttributeView, attributeViewEntity] = (*attributeViewMapper)(nil)
