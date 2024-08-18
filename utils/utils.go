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
		"id":        user.ID,
		"first_name": user.FirstName,
		"last_name":  user.LastName,
		"email":     user.Email,
		"createdAt": user.CreatedAt,
	}
}
