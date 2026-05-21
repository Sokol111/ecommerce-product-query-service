package productview

// FacetValue represents a single value within an attribute facet with its product count
type FacetValue struct {
	Slug  string // Option slug (for single/multiple types)
	Value string // Raw value as string
	Count int    // Number of products with this value
}

// AttributeFacet represents facet data for a single attribute
type AttributeFacet struct {
	AttributeID string       // Attribute definition ID
	Slug        string       // Attribute slug for filtering
	Values      []FacetValue // Available values with product counts
}

// PriceRange represents the min/max price boundaries
type PriceRange struct {
	Min float64
	Max float64
}

// FacetsResult holds the complete facets computation result
type FacetsResult struct {
	Facets     []AttributeFacet
	PriceRange PriceRange
}
