package product

import (
	"encoding/json"
	"net/http"

	"github.com/Hyp9r/csv-processing-service/port/rest/common"
)

func (ctrl *Controller) ListProductsHandler(w http.ResponseWriter, r *http.Request) {
	ctrl.logger.Info().Msg("inside of list products handler")
	encoder := json.NewEncoder(w)
	products, err := ctrl.productService.List()
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
	err = encoder.Encode(fromDomainToArray(products))
	if err != nil {
		ctrl.logger.Err(err).Msg("error while serializing response")
		return
	}
	w.WriteHeader(http.StatusOK)
}