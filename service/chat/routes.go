package chat

import (
	"ChiragKr04/go-backend/service"
	"net/http"

	"github.com/gorilla/mux"
)

func (h *Handler) ChatRoutes(router *mux.Router) {
	router.HandleFunc("/chat/history/{roomId}", service.WrapWithAuth(h.handleGetChatHistory)).Methods(http.MethodGet, http.MethodOptions)
}
