package Middlewares

import (
	"errors"
	"fmt"
	"net/http"

	"216/internal/orchestrator/Response"
	"216/internal/orchestrator/Services"
)

func SetMiddlewareJSON(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next(w, r)
	}
}

func SetMiddlewareAuthentication(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := Services.TokenValid(r)
		if err != nil {
			fmt.Println(err)
			Response.ErrorResponse(w, http.StatusUnauthorized, errors.New("Unauthorized"))
			return
		}
		next(w, r)
	}
}
