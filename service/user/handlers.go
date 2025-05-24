package user

import (
	"ChiragKr04/go-backend/service/auth"
	"ChiragKr04/go-backend/types"
	"ChiragKr04/go-backend/utils"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-playground/validator"
	"github.com/gorilla/mux"
)

func (h *Handler) handleLogin(w http.ResponseWriter, r *http.Request) {
	// Handle login
	// get email and password from request body
	var payload types.UserLoginRequest
	if err := utils.ParseJSON(r, &payload); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("validation error: %s", errors))
		return
	}
	if err := utils.Validate.Struct(payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}
	user, err := h.repo.GetUserByEmail(payload.Email)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, errors.New("user not found"))
		return
	}
	if !auth.CheckPasswordHash(payload.Password, user.Password) {
		utils.WriteError(w, http.StatusBadRequest, errors.New("invalid password"))
		return
	}
	secret := []byte("secret")
	token, err := auth.CreateJWTToken(secret, user.ID)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, errors.New("error creating token"))
		return
	}
	// remove password from user data
	user.Password = ""
	response := map[string]interface{}{
		"user":    utils.ReturnUserWithoutPassword(*user),
		"token":   token,
		"message": "success",
	}
	utils.WriteJSON(w, http.StatusOK, response)
}

func (h *Handler) handleRegister(w http.ResponseWriter, r *http.Request) {
	var payload types.UserRegisterRequest
	if err := utils.ParseJSON(r, &payload); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("validation error: %s", errors))
		return
	}
	if err := utils.Validate.Struct(payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}
	_, err := h.repo.GetUserByEmail(payload.Email)
	if err == nil {
		utils.WriteError(w, http.StatusBadRequest, errors.New("user already exists"))
		return
	}
	hashPassword, err := auth.HashPassword(payload.Password)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	err = h.repo.CreateUser(types.User{
		FirstName: payload.FirstName,
		LastName:  payload.LastName,
		Email:     payload.Email,
		Password:  hashPassword,
		Username:  payload.Username,
	})
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	utils.WriteJSON(w, http.StatusCreated, map[string]string{"message": "user created"})
}

func (h *Handler) handleGetProfile(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID, err := strconv.Atoi(vars["userId"])
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, errors.New("invalid user id"))
		return
	}
	// Retrieve user profile
	user, err := h.repo.GetUserByID(userID)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Return user profile
	utils.WriteJSON(w, http.StatusOK, returnUserWithToken(*user))
}

func (h *Handler) handleUpdateProfile(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID, err := strconv.Atoi(vars["userId"])
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, errors.New("invalid user id"))
		return
	}
	var payload types.UserUpdateRequest
	if err := utils.ParseJSON(r, &payload); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("validation error: %s", errors))
		return
	}
	if err := utils.Validate.Struct(payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}
	user, err := h.repo.GetUserByID(userID)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, errors.New("user not found"))
		return
	}
	if payload.FirstName != "" {
		user.FirstName = payload.FirstName
	}
	if payload.LastName != "" {
		user.LastName = payload.LastName
	}
	if payload.Email != "" {
		user.Email = payload.Email
	}
	user, err = h.repo.UpdateUser(*user)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, returnUserWithToken(*user))
}

func returnUserWithToken(user types.User) map[string]interface{} {
	secret := []byte("secret")
	token, err := auth.CreateJWTToken(secret, user.ID)
	if err != nil {
		return nil
	}
	return map[string]interface{}{
		"message": "success",
		"token":   token,
		"user":    utils.ReturnUserWithoutPassword(user),
	}
}
