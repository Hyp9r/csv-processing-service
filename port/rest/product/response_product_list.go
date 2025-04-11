package product

import "github.com/Hyp9r/csv-processing-service/domain/product"

type ListResponse struct {
	Count int `json:"count"`
	Products []ProductResponse `json:"products"`
}

func fromDomainToArray(products []product.Product) ListResponse {
	var ps []ProductResponse
	for _, p := range products {
		ps = append(ps, fromDomain(p))
	}
	return ListResponse{
		Count: len(ps),
		Products: ps,
	}
}