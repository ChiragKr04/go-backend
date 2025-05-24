package rooms

import (
	"ChiragKr04/go-backend/service/auth"
	"ChiragKr04/go-backend/types"
	"ChiragKr04/go-backend/utils"
	"errors"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func (h *Handler) handleCreateRoom(w http.ResponseWriter, r *http.Request) {
	user, err := h.getUserFromContext(w, r)
	if err != nil {
		return
	}
	var payload types.RoomCreateRequest
	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid request format: %v", err))
		return
	}
	if err := utils.Validate.Struct(payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}
	id, err := h.repo.CreateRoom(user, payload)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	//  convert id to int64
	id64 := int64(id)
	room, err := h.repo.GetRoomById(id64)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, room)
}

func (h *Handler) HandleGetRoomById(w http.ResponseWriter, r *http.Request) {
	_, err := h.getUserFromContext(w, r)
	if err != nil {
		return
	}
	roomId := mux.Vars(r)["roomId"]
	room, err := h.repo.GetRoomByRoomId(roomId)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, room)
}

func (h *Handler) getUserFromContext(w http.ResponseWriter, r *http.Request) (*types.User, error) {
	userIDValue := r.Context().Value(auth.UserIDKey)
	if userIDValue == nil {
		utils.WriteError(w, http.StatusUnauthorized, errors.New("unauthorized"))
		return nil, nil
	}

	// Convert userIDValue to int
	userID, ok := userIDValue.(int)
	if !ok {
		utils.WriteError(w, http.StatusInternalServerError, errors.New("error converting userID to int"))
		return nil, nil
	}
	user, err := h.userRepo.GetUserByID(userID)
	if err != nil {
		utils.WriteError(w, http.StatusNotFound, errors.New("user not found"))
		return nil, nil
	}
	return user, nil
}
