package product

import "github.com/Hyp9r/csv-processing-service/domain/product"

type ImportResponse struct {
	InsertedRows int `json:"inserted"`
	InvalidRows int `json:"invalid"`
	Errors []product.CSVFileRowError `json:"errors"`
}