package service

import (
	"ChiragKr04/go-backend/service/auth"
	"net/http"
)

func WrapWithAuth(handlerFunc http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		auth.AuthMiddleware()(handlerFunc).ServeHTTP(w, r)
	}
}