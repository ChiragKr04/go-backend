package websocket

import (
	"ChiragKr04/go-backend/service/user"
	"ChiragKr04/go-backend/types"
	"net/http"

	"github.com/gorilla/mux"
)

type Handler struct {
	Hub *HubFile

	UserRepo *user.UserRepository
}

func NewHandler(repo *user.UserRepository) *Handler {
	return &Handler{
		Hub:      newHub(),
		UserRepo: repo,
	}
}

type HubType struct {
	Hub *types.Hub
}

func newHub() *HubFile {
	return &HubFile{
		HubType: &types.Hub{
			Broadcast:  make(chan []byte),
			Register:   make(chan *types.Client),
			Unregister: make(chan *types.Client),
			Clients:    make(map[*types.Client]bool),
		},
	}
}

func (h *Handler) WebsocketRoutes(router *mux.Router) {
	hub := newHub()
	go hub.Run()

	router.HandleFunc("/ws/{userId}", func(w http.ResponseWriter, r *http.Request) {
		h.serveWs(hub, w, r)
	})
}
