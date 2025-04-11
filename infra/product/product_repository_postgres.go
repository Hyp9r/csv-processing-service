package postgres

import (
	"database/sql"
	"strings"

	"github.com/Hyp9r/csv-processing-service/domain/product"
	"github.com/rs/zerolog"
)


type ProductRepository struct {
	*sql.DB
	logger *zerolog.Logger
}

// Actual table definition
type Product struct {
	ID string
	Name string
	Category string
	Price float32
	Description string
	BrandName string
	StockQuantity int
	Manufacturer string
	Sku string
	Weight float32
	Color string
}

func NewProductRepository(db *sql.DB, logger *zerolog.Logger) *ProductRepository {
	return &ProductRepository{
		db,
		logger,
	}
}

var _ product.ProductRepository = (*ProductRepository)(nil)

func (pr *ProductRepository) Create(product product.Product) error {
	query := "INSERT INTO \"products\"(name, category, price, description, brand_name, stock_quantity, manufacturer, sku, weight, color) VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)"
	err := pr.DB.QueryRow(query, product.Name, product.Category, product.Price, product.Description, product.BrandName, product.StockQuantity, product.Manufacturer, product.Sku, product.Weight, product.Color).Err()
	if err != nil {
		return err
	}
	return nil
}

func (pr *ProductRepository) Get(productID string) (product.Product, error) {
	query := "SELECT * FROM \"products\" WHERE id = $1"
	row := pr.DB.QueryRow(query, productID)

	var p Product

	err := row.Scan(&p.ID, &p.Name, &p.Category, &p.Price, &p.Description, &p.BrandName, &p.StockQuantity, &p.Manufacturer, &p.Sku, &p.Weight, &p.Color)
	if err != nil {
		if err == sql.ErrNoRows {
			pr.logger.Err(err).Msgf("error, no rows in database with id: %s", productID)
			return product.Product{}, err
		}
		pr.logger.Err(err).Msg("error while trying to scan product details")
		return product.Product{}, err
	}
	return p.toDomain(), nil
}

func (pr *ProductRepository) List() ([]product.Product, error) {
	query := "SELECT * FROM \"products\""
	rows, err := pr.DB.Query(query)
	if err != nil {
		pr.logger.Err(err).Msg("error while trying to query list of all products")
		return nil, err
	}
	defer rows.Close()

	var products []product.Product

	for rows.Next() {
		var p Product
		err := rows.Scan(&p.ID, &p.Name, &p.Category, &p.Price, &p.Description, &p.BrandName, &p.StockQuantity, &p.Manufacturer, &p.Sku, &p.Weight, &p.Color)
		if err != nil {
			pr.logger.Err(err).Msg("error while scanning product rows")
			return nil, err
		}
		products = append(products, p.toDomain())
	}

	if err := rows.Err(); err != nil {
		pr.logger.Err(err).Msg("error while iterating through product rows")
		return nil, err
	}

	return products, nil
}

func (pr *ProductRepository) Update(product product.Product) error {
	query := "UPDATE \"products\" SET name = COALESCE($1, name), category = COALESCE($2, category), price = COALESCE($3, price), description = COALESCE($4, description), brand_name = COALESCE($5, brand_name), stock_quantity = COALESCE($6, stock_quantity) , manufacturer = COALESCE($7, manufacturer) , sku = COALESCE($8, sku) , weight = COALESCE($9, weight) , color = COALESCE($10, color) WHERE id = $11"
	dbResp := pr.DB.QueryRow(query, &product.Name, &product.Category, &product.Price, &product.Description, &product.BrandName, &product.StockQuantity, &product.Sku, &product.Weight, &product.Color)
	if dbResp.Err() != nil {
		if dbResp.Err() == sql.ErrNoRows{
			pr.logger.Err(dbResp.Err()).Msgf("product with id: %s doesn't exist", product.ID)
			return dbResp.Err()
		} else {
			pr.logger.Err(dbResp.Err()).Msg("failed to edit product")
			return dbResp.Err()
		}
	}
	return nil
}

func (pr *ProductRepository) Delete(productID string) error {
	query := "DELETE FROM \"products\" WHERE id = $1"
	err := pr.DB.QueryRow(query, productID)
	if err != nil {
		return err.Err()
	}
	return nil
}

func (pr *ProductRepository) BatchInsert(queryFragments []string, insertValues []interface{}) error {
	query := `INSERT INTO products (product_name, product_category, product_price, 
	product_description, brand_name, stock_quantity, manufacturer, sku, weight, color)
	VALUES ` + strings.Join(queryFragments, ",") + `ON CONFLICT(sku) DO NOTHING`

	_, err := pr.DB.Exec(query, insertValues...)
	if err != nil {
		pr.logger.Err(err).Msg("batch insert failed")
		return err
	}
	return nil
}

func (p *Product) toDomain() product.Product {
	return product.Product{
		ID: p.ID,
		Name: p.Name,
		Category: p.Category,
		Price: p.Price,
		Description: p.Description,
		BrandName: p.BrandName,
		StockQuantity: p.StockQuantity,
		Manufacturer: p.Manufacturer,
		Sku: p.Sku,
		Weight: p.Weight,
		Color: p.Color,
	}
}