package rooms

import (
	"ChiragKr04/go-backend/service"
	"ChiragKr04/go-backend/types"
	"net/http"

	"github.com/gorilla/mux"
)

type Handler struct {
	repo types.RoomRepository
}

func NewHandler(repo types.RoomRepository) *Handler {
	return &Handler{
		repo: repo,
	}
}

func (h *Handler) RoomRoutes(router *mux.Router) {
	router.HandleFunc("/create-room", service.WrapWithAuth(h.handleCreateRoom)).Methods(http.MethodPost)
}
