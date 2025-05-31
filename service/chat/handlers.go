package chat

import (
	"ChiragKr04/go-backend/types"
	"ChiragKr04/go-backend/utils"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type Handler struct {
	chatRepo types.ChatRepository
}

func NewHandler(chatRepo types.ChatRepository) *Handler {
	return &Handler{
		chatRepo: chatRepo,
	}
}

func (h *Handler) handleGetChatHistory(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	roomId := vars["roomId"]

	// Get pagination parameters
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	limit := 50 // default limit
	offset := 0 // default offset

	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
			limit = l
		}
	}

	if offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil && o >= 0 {
			offset = o
		}
	}

	chats, err := h.chatRepo.GetChatsByRoomId(roomId, limit, offset)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"chats":  chats,
		"limit":  limit,
		"offset": offset,
	})
}
