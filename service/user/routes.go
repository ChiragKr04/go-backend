package user

import (
	"ChiragKr04/go-backend/service/auth"
	"ChiragKr04/go-backend/types"
	"net/http"

	"github.com/gorilla/mux"
)

type Handler struct {
	repo types.UserRepository
}

func NewHandler(repo types.UserRepository) *Handler {
	return &Handler{
		repo: repo,
	}
}

func (h *Handler) UserRoutes(router *mux.Router) {
	router.HandleFunc("/login", h.handleLogin).Methods(http.MethodPost)
	router.HandleFunc("/register", h.handleRegister).Methods(http.MethodPost)
	router.HandleFunc("/get-profile/{userId}", func(w http.ResponseWriter, r *http.Request) {
		auth.AuthMiddleware()(http.HandlerFunc(h.handleGetProfile)).ServeHTTP(w, r)
	}).Methods(http.MethodGet)
	router.HandleFunc("/update-profile/{userId}", func(w http.ResponseWriter, r *http.Request) {
		auth.AuthMiddleware()(http.HandlerFunc(h.handleUpdateProfile)).ServeHTTP(w, r)
	}).Methods(http.MethodPut)
}
