package rooms

import "net/http"


func (h *Handler) handleCreateRoom(w http.ResponseWriter, r *http.Request) {
	h.repo.CreateRoom()
}
