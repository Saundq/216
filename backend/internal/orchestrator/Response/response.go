package Response

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func JsonResponse(w http.ResponseWriter, statusCode int, data interface{}) {
	w.WriteHeader(statusCode)
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		fmt.Fprintf(w, "%s", err.Error())
	}
}

func ErrorResponse(w http.ResponseWriter, statusCode int, err error) {
	if err != nil {
		JsonResponse(w, statusCode, struct {
			Error string `json:"error"`
		}{
			Error: err.Error(),
		})
		return
	}
	JsonResponse(w, http.StatusBadRequest, nil)
}
