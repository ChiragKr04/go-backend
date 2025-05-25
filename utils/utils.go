package utils

import (
	"ChiragKr04/go-backend/types"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-playground/validator"
	"github.com/golang-jwt/jwt"
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

func ReturnUserWithoutPassword(user types.User) types.User {
	return types.User{
		ID:        user.ID,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		Username:  user.Username,
	}
}

func GetUserFromContext(w http.ResponseWriter, r *http.Request, userRepo types.UserRepository) (*types.User, error) {
	userIDValue := r.Context().Value(types.UserIDKey)
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

func GetUserIDFromToken(w http.ResponseWriter, r *http.Request) (int, error) {
	tokenValue := r.URL.Query().Get("token")
	if tokenValue == "" {
		WriteError(w, http.StatusBadRequest, errors.New("token is required"))
		return 0, nil
	}
	// get user id from token
	tokenValue = strings.TrimPrefix(tokenValue, "Bearer ")

	// Parse and validate the token
	token, err := jwt.Parse(tokenValue, func(token *jwt.Token) (interface{}, error) {
		// Check if the signing method is correct
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte("secret"), nil
	})

	if err != nil {
		return 0, err
	}

	// Extract the claims (payload)
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// Extract the user ID from the claims - use "userId" key directly
		if userIDStr, ok := claims["userId"].(string); ok {
			userID, err := strconv.Atoi(userIDStr)
			if err != nil {
				return 0, errors.New("invalid user ID format")
			}
			return userID, nil
		}
		return 0, errors.New("invalid token claims")
	}

	return 0, errors.New("invalid token")

}
