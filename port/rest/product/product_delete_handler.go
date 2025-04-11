package product

import (
	"encoding/json"
	"net/http"

	"github.com/Hyp9r/csv-processing-service/port/rest/common"
)

func (ctrl *Controller) DeleteProductHandler(w http.ResponseWriter, r *http.Request) {
	ctrl.logger.Info().Msg("inside of delete product handler")
	productID := r.PathValue("id")
	encoder := json.NewEncoder(w)
	err := ctrl.productService.Delete(productID)
	if err != nil {
		if err != nil {
			resp := common.NewAPIError(APIPath, DELETE, err.Error())
			encodeErr := encoder.Encode(resp)
			if encodeErr != nil {
				ctrl.logger.Err(err).Msg("error while serializing error response")
				return
			}
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
	w.WriteHeader(http.StatusOK)
	ctrl.logger.Info().Msg("successfully deleted product")
}