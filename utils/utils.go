package utils

import (
	"ChiragKr04/go-backend/types"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-playground/validator"
)

var Validate = validator.New()

func ParseJSON(r *http.Request, payload any) error {
	if r.Body == nil {
		return errors.New("please send a request body")
	}
	return json.NewDecoder(r.Body).Decode(payload)
}

func WriteJSON(w http.ResponseWriter, status int, output any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(output)
}

func WriteError(w http.ResponseWriter, status int, err error) {
	WriteJSON(w, status, map[string]string{"error": err.Error()})
}

func ReturnUserWithoutPassword(user types.User) map[string]interface{} {
	return map[string]interface{}{
		"id":         user.ID,
		"first_name": user.FirstName,
		"last_name":  user.LastName,
		"email":      user.Email,
		"createdAt":  user.CreatedAt,
		"username":   user.Username,
	}
}

func GetUserFromContext(w http.ResponseWriter, r *http.Request, userRepo types.UserRepository) (*types.User, error) {
	userIDValue := r.Context().Value("userId")
	if userIDValue == nil {
		WriteError(w, http.StatusUnauthorized, errors.New("unauthorized"))
		return nil, nil
	}

	// Convert userIDValue to int
	userID, ok := userIDValue.(int)
	if !ok {
		WriteError(w, http.StatusInternalServerError, errors.New("error converting userID to int"))
		return nil, nil
	}
	user, err := userRepo.GetUserByID(userID)
	if err != nil {
		WriteError(w, http.StatusNotFound, errors.New("user not found"))
		return nil, nil
	}
	return user, nil
}
