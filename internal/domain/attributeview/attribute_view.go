package attributeview

import "time"

// AttributeType represents the type of attribute values
type AttributeType string

const (
	AttributeTypeSingle   AttributeType = "single"
	AttributeTypeMultiple AttributeType = "multiple"
	AttributeTypeRange    AttributeType = "range"
	AttributeTypeBoolean  AttributeType = "boolean"
	AttributeTypeText     AttributeType = "text"
)

// AttributeOption represents an option for single/multiple type attributes
type AttributeOption struct {
	Slug      string  // URL-friendly identifier (immutable)
	Name      string  // Display name
	ColorCode *string // Hex color code for color options
	SortOrder int     // Sort order for display
}

// AttributeView - read model for attribute master data (CQRS query side)
// Used to store mutable attribute data that can be joined with product attributes at read time
type AttributeView struct {
	ID         string            // Unique attribute identifier (UUID)
	Version    int               // Entity version for conflict resolution
	Slug       string            // URL-friendly identifier (immutable)
	Name       string            // Display name of the attribute
	Type       AttributeType     // Attribute type (single, multiple, range, boolean, text)
	Unit       *string           // Unit of measurement for range type attributes
	Enabled    bool              // Whether the attribute is enabled for display
	ModifiedAt time.Time         // Last modification timestamp
	Options    []AttributeOption // Available options for single/multiple type attributes
}

// Reconstruct creates an AttributeView from persistence data
func Reconstruct(
	id string,
	version int,
	slug string,
	name string,
	attrType AttributeType,
	unit *string,
	enabled bool,
	modifiedAt time.Time,
	options []AttributeOption,
) *AttributeView {
	return &AttributeView{
		ID:         id,
		Version:    version,
		Slug:       slug,
		Name:       name,
		Type:       attrType,
		Unit:       unit,
		Enabled:    enabled,
		ModifiedAt: modifiedAt,
		Options:    options,
	}
}
