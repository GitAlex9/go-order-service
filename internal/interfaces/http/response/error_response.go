package response

import (
	"encoding/json"
	"net/http"
)

type ErrorResponse struct {
	Error  string      `json:"error"`
	Fields []FieldItem `json:"fields,omitempty"`
}

type FieldItem struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func JSONError(w http.ResponseWriter, status int, message string, fields []FieldItem) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(ErrorResponse{Error: message, Fields: fields})
}
