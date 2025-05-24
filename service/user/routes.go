package user

import (
	"ChiragKr04/go-backend/service"
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
	router.HandleFunc("/login", h.handleLogin).Methods(http.MethodPost, http.MethodOptions)
	router.HandleFunc("/register", h.handleRegister).Methods(http.MethodPost, http.MethodOptions)
	router.HandleFunc(
		"/get-profile/{userId}",
		service.WrapWithAuth(h.handleGetProfile)).Methods(http.MethodGet, http.MethodOptions)
	router.HandleFunc(
		"/update-profile/{userId}",
		service.WrapWithAuth(h.handleUpdateProfile)).Methods(http.MethodPut, http.MethodOptions)
}
