package rooms

import (
	"ChiragKr04/go-backend/service"
	"ChiragKr04/go-backend/types"
	"net/http"

	"github.com/gorilla/mux"
)

type Handler struct {
	repo     types.RoomRepository
	userRepo types.UserRepository
}

func NewHandler(repo types.RoomRepository, userRepo types.UserRepository) *Handler {
	return &Handler{
		repo:     repo,
		userRepo: userRepo,
	}
}

func (h *Handler) RoomRoutes(router *mux.Router) {
	router.HandleFunc("/create-room", service.WrapWithAuth(h.handleCreateRoom)).Methods(http.MethodPost)
	router.HandleFunc("/get-room-by-roomid/{roomId}", service.WrapWithAuth(h.HandleGetRoomById)).Methods(http.MethodGet)
}
