package product

import (
	"encoding/csv"
	"fmt"
	"io"
	"mime/multipart"
	"runtime"
	"strconv"
	"sync"

	"github.com/rs/zerolog"
)

type ProductService struct {
	repository ProductRepository
	logger *zerolog.Logger
}

func NewProductService(repository ProductRepository, logger *zerolog.Logger) *ProductService {
	return &ProductService{
		repository: repository,
		logger: logger,
	}
}

func (ps *ProductService) Get(productID string) (Product, error) {
	product, err := ps.repository.Get(productID)
	if err != nil {
		ps.logger.Err(err).Str("product_id", productID).Msg("error while trying to query product")
		return Product{}, err
	}
	return product, nil
}

func (ps *ProductService) Delete(productID string) error {
	err := ps.repository.Delete(productID)
	if err != nil {
		ps.logger.Err(err).Str("product_id", productID).Msg("production deletion failed")
		return err
	}
	return nil
}

func (ps *ProductService) List() ([]Product, error) {
	products, err := ps.repository.List()
	if err != nil {
		ps.logger.Err(err).Msg("failed to fetch all products")
		return nil, err
	}
	return products, nil
}

func (ps *ProductService) Update(product Product) error {
	err := ps.repository.Update(product)
	if err != nil {
		ps.logger.Err(err).Str("product_id", product.ID).Msg("failed to update product")
		return err
	}
	return nil
}

func (ps *ProductService) BatchInsert(queryFragments []string, insertValues []interface{}) error {
	err := ps.repository.BatchInsert(queryFragments, insertValues)
	if err != nil {
		ps.logger.Err(err).Msg("error while trying to batch insert products")
		return err
	}
	return nil
}

func (ps *ProductService) ProcessProductImport(file multipart.File) (ProcessingResult, error) {
	reader := csv.NewReader(file)
	headers, err := reader.Read()
	if err != nil {
		ps.logger.Err(err).Msg("error reading CSV header")
		return ProcessingResult{}, err
	}
	ps.logger.Info().Msgf("headers: %v\n", headers)

	// prepare channels
	rowChan := make(chan IndexedRow)
	resultChan := make(chan CSVFileRow)
	errorChan := make(chan CSVFileRowError)

	var wg sync.WaitGroup
	numberOfWorkers := runtime.NumCPU()

	for i := 0; i < numberOfWorkers; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			for row := range rowChan {
				line := row.LineNumber
				hasErrors := false
				if len(row.Row) < 10 {
					ps.logger.Err(err).Msg("invalid row skipping...")
					continue
				}

				price, err := strconv.ParseFloat(row.Row[2], 64)
				if err != nil {
					fmt.Printf("price couldn't be parsed %v\n", err)
					rowError := CSVFileRowError {
						LineNumber: line,
						Error: "price has wrong format, can't be parsed",
						Column: "Price",
						Value: row.Row[2],
					}
					errorChan <- rowError
					hasErrors = true
				}
				stock, err := strconv.Atoi(row.Row[5])
				if err != nil {
					fmt.Printf("stock couldn't be parsed %v\n", err)
					rowError := CSVFileRowError {
						LineNumber: line,
						Error: "stock has wrong format, can't be parsed",
						Column: "Stock",
						Value: row.Row[5],
					}
					errorChan <- rowError
					hasErrors = true
				}
				weight, err := strconv.ParseFloat(row.Row[8], 64)
				if err != nil {
					fmt.Printf("weight couldn't be parsed %v\n", err)
					rowError := CSVFileRowError {
						LineNumber: line,
						Error: "weight has wrong format, can't be parsed",
						Column: "Weight",
						Value: row.Row[8],
					}
					errorChan <- rowError
					hasErrors = true
				}

				if hasErrors {
					continue
				}

				csvProcessedRow := CSVFileRow {
					ProductName: row.Row[0],
					ProductCategory: row.Row[1],
					ProductPrice: price,
					ProductDescription: row.Row[3],
					BrandName: row.Row[4],
					StockQuantity: stock,
					Manufacturer: row.Row[6],
					Sku: row.Row[7],
					Weight: weight,
					Color: row.Row[9],
				}
				resultChan <- csvProcessedRow
			}
		}(i)
	}

	// feed rows into channel that's going to be consumed by workers
	go func() {
		counter := 2
		for {
			record, err := reader.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				continue // skip bad row
			}
			rowChan <- IndexedRow{ LineNumber: counter, Row: record}
			counter++
		}
		close(rowChan)
	}()

	go func() {
		wg.Wait()
		close(resultChan)
		close(errorChan)
	}()

	var wgCollect sync.WaitGroup
	var csvRows []CSVFileRow
	var csvErrors []CSVFileRowError
	skuMap := make(map[string]struct{})

	wgCollect.Add(2)

	go func() {
		defer wgCollect.Done()

		valueStrings := []string{}
		var insertValues []interface{}
		counter := 0


		for row := range resultChan {
			_, exists := skuMap[row.Sku]
			if !exists {
				csvRows = append(csvRows, row)
				skuMap[row.Sku] = struct{}{}
				valueString := fmt.Sprintf(
					"($%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d)",
					counter*10+1, counter*10+2, counter*10+3,counter*10+4,
					counter*10+5,counter*10+6,counter*10+7,counter*10+8,
					counter*10+9, counter*10+10,
				)
				valueStrings = append(valueStrings, valueString)

				insertValues = append(
					insertValues, row.ProductName, row.ProductCategory, row.ProductPrice,
					row.ProductDescription, row.BrandName,row.StockQuantity,
					row.Manufacturer, row.Sku, row.Weight, row.Color,
				)
				counter++
			}
		}
		
		if counter == 0 {
			return
		}

		err := ps.repository.BatchInsert(valueStrings, insertValues)
		if err != nil {
			ps.logger.Err(err).Msg("failed to batch insert processed csv file results")
			return
		}
	}()

	go func() {
		defer wgCollect.Done()
		for err := range errorChan {
			csvErrors = append(csvErrors, err)
		}
	}()

	wgCollect.Wait()

	resp := ProcessingResult{
		ImportedRows: len(csvRows),
		Errors: csvErrors,
		InvalidRows: len(csvErrors),
	}

	return resp, nil
}