package websocket

import (
	"ChiragKr04/go-backend/service/rooms"
	"ChiragKr04/go-backend/service/user"
	"ChiragKr04/go-backend/types"
	"net/http"

	"github.com/gorilla/mux"
)

type Handler struct {
	HubManager *HubManager
	UserRepo   *user.UserRepository
	RoomRepo   *rooms.RoomsRepository
	ChatRepo   types.ChatRepository
	WebrtcRepo types.WebrtcRepository
}

func NewHandler(
	repo *user.UserRepository,
	roomRepo *rooms.RoomsRepository,
	chatRepo types.ChatRepository,
	webrtcRepo types.WebrtcRepository) *Handler {
	return &Handler{
		HubManager: NewHubManager(),
		UserRepo:   repo,
		RoomRepo:   roomRepo,
		ChatRepo:   chatRepo,
		WebrtcRepo: webrtcRepo,
	}
}

func (h *Handler) WebsocketRoutes(router *mux.Router) {
	router.HandleFunc("/ws/{roomId}", h.serveWs).Methods(http.MethodGet, http.MethodOptions)
}
