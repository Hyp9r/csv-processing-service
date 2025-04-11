package product

import (
	"encoding/json"
	"net/http"
)

func (ctrl *Controller) ImportProductsHandler(w http.ResponseWriter, r *http.Request) {
	file, _, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Error retrieving file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	result, err := ctrl.productService.ProcessProductImport(file)
	if err != nil {
		ctrl.logger.Err(err).Msg("failed to process file")
		w.WriteHeader(http.StatusBadRequest)
	}

	w.WriteHeader(http.StatusCreated)

	resp := ImportResponse{
		Errors: result.Errors,
		InsertedRows: result.ImportedRows,
		InvalidRows: result.InvalidRows,
	}
	json.NewEncoder(w).Encode(resp)
	ctrl.logger.Info().Msgf("Successfully processed: %d rows\n", result.ImportedRows)
}