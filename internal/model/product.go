package model

var NotFoundError = errors.New("not found")

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

type Service interface {

	// can return NotFoundError
	GetById(ctx context.Context, id string) (*Product, error)

	GetAll(ctx context.Context) ([]*Product, error)
}
