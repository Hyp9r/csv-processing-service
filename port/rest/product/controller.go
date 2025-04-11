package product

import (
	"net/http"

	"github.com/Hyp9r/csv-processing-service/domain/product"
	"github.com/rs/zerolog"
)

type Controller struct {
	productService *product.ProductService
	logger      *zerolog.Logger
	router *http.ServeMux
}

func RegisterRoutes(productService *product.ProductService, logger *zerolog.Logger, router *http.ServeMux) {
	ctrl := &Controller{
		productService: productService,
		logger: logger,
		router: router,
	}
	ctrl.router.HandleFunc("GET /products", ctrl.ListProductsHandler)
	ctrl.router.HandleFunc("GET /products/{id}", ctrl.GetProductHandler)
	ctrl.router.HandleFunc("DELETE /products/{id}", ctrl.DeleteProductHandler)
	ctrl.router.HandleFunc("POST /products/import", ctrl.ImportProductsHandler)
}