package productdetail

type entity struct {
	ID       string `bson:"_id"`
	Version  int
	Name     string
	Price    float32
	Quantity int
	Enabled  bool
}
