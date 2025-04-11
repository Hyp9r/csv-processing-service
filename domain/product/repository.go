package product

type ProductRepository interface {
	Get(productID string) (Product, error)
	Update(product Product) error
	Delete(productID string) error
	List() ([]Product, error)
	Create(product Product) error
	BatchInsert(queryFragments []string, insertValues []interface{}) error
}