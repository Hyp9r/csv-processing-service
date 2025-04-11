package product

type IndexedRow struct {
	LineNumber int
	Row []string
}

type CSVFileRow struct {
	ProductName string
	ProductCategory string
	ProductPrice float64
	ProductDescription string
	BrandName string
	StockQuantity int
	Manufacturer string
	Sku string
	Weight float64
	Color string
}

type CSVFileRowError struct {
	LineNumber int
	Error string
	Column string
	Value string
}

type ProcessingResult struct {
	ImportedRows int
	InvalidRows int
	Errors []CSVFileRowError
}