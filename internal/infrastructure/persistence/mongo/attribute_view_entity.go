package mongo

import "time"

// attributeOptionEntity represents the MongoDB subdocument structure for attribute options
type attributeOptionEntity struct {
	Slug      string  `bson:"slug"`
	Name      string  `bson:"name"`
	ColorCode *string `bson:"colorCode,omitempty"`
	SortOrder int     `bson:"sortOrder"`
}

// attributeViewEntity represents the MongoDB document structure for attribute master data
type attributeViewEntity struct {
	ID         string                  `bson:"_id"`
	Version    int                     `bson:"version"`
	Slug       string                  `bson:"slug"`
	Name       string                  `bson:"name"`
	Type       string                  `bson:"type"`
	Unit       *string                 `bson:"unit,omitempty"`
	Enabled    bool                    `bson:"enabled"`
	ModifiedAt time.Time               `bson:"modifiedAt"`
	Options    []attributeOptionEntity `bson:"options,omitempty"`
}
