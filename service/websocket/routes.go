package websocket

import (
	"ChiragKr04/go-backend/service/rooms"
	"ChiragKr04/go-backend/service/user"
	"net/http"

	"github.com/gorilla/mux"
)

type Handler struct {
	HubManager *HubManager
	UserRepo   *user.UserRepository
	RoomRepo   *rooms.RoomsRepository
}

func NewHandler(
	repo *user.UserRepository,
	roomRepo *rooms.RoomsRepository) *Handler {
	return &Handler{
		HubManager: NewHubManager(),
		UserRepo:   repo,
		RoomRepo:   roomRepo,
	}
}

func (h *Handler) WebsocketRoutes(router *mux.Router) {
	router.HandleFunc("/ws/{roomId}", h.serveWs).Methods(http.MethodGet, http.MethodOptions)
}
