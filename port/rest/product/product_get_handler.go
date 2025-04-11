package product

import (
	"encoding/json"
	"net/http"

	"github.com/Hyp9r/csv-processing-service/port/rest/common"
)

func (ctrl *Controller) GetProductHandler(w http.ResponseWriter, r *http.Request) {
	ctrl.logger.Info().Msg("inside of get product handler")
	productID := r.PathValue("id")
	encoder := json.NewEncoder(w)
	product, err := ctrl.productService.Get(productID)
	if err != nil {
		resp := common.NewAPIError(APIPath, GET, err.Error())
		encodeErr := encoder.Encode(resp)
		if encodeErr != nil {
			ctrl.logger.Err(err).Msg("error while serializing error response")
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = encoder.Encode(fromDomain(product))
	if err != nil {
		ctrl.logger.Err(err).Msg("error while serializing response")
		return
	}
	w.WriteHeader(http.StatusOK)
}