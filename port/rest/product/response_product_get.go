package product

import "github.com/Hyp9r/csv-processing-service/domain/product"

type ProductResponse struct {
	ID string `json:"id"`
	Name string `json:"name"`
	Category string `json:"category"`
	Price float32 `json:"price"`
	Description string `json:"description"`
	BrandName string `json:"brandName"`
	StockQuantity int `json:"stockQuantity"`
	Manufacturer string `json:"manufacturer"`
	Sku string `json:"sku"`
	Weight float32 `json:"weight"`
	Color string `json:"color"`
}

func fromDomain(product product.Product) ProductResponse {
	return ProductResponse{
		ID: product.ID,
		Name: product.Name,
		Category: product.Category,
		Price: product.Price,
		Description: product.Description,
		BrandName: product.BrandName,
		StockQuantity: product.StockQuantity,
		Manufacturer: product.Manufacturer,
		Sku: product.Sku,
		Weight: product.Weight,
		Color: product.Color,
	}
}

