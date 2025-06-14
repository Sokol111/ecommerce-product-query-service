package product

import "time"

type Product struct {
	ID         string
	Version    int
	Name       string
	Price      float32
	Quantity   int
	Enabled    bool
	CreatedAt  time.Time
	ModifiedAt time.Time
}
