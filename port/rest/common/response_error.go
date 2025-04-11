package common

type ErrorResponse struct {
	Path string `json:"path"`
	Operation string `json:"operation"`
	Message string `json:"message"`
}

func NewAPIError(path string, operation string, message string) *ErrorResponse {
	return &ErrorResponse{Path: path, Operation: operation, Message: message}
}