package auth

import (
	"ChiragKr04/go-backend/utils"
	"context"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/golang-jwt/jwt"
)

type contextKey string

const userIDKey contextKey = "userId"

func AuthMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
				utils.WriteError(w, http.StatusUnauthorized, errors.New("missing or invalid authorization header"))
				return
			}

			token := strings.TrimPrefix(authHeader, "Bearer ")

			// Validate the token and extract the user ID
			secret := []byte("secret")
			userID, err := validateJWT(token, secret)
			if err != nil {
				utils.WriteError(w, http.StatusUnauthorized, err)
				return
			}

			// Add the user ID to the context

			ctx := context.WithValue(r.Context(), userIDKey, userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func validateJWT(tokenString string, secret []byte) (int, error) {
	// Parse and validate the token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Check if the signing method is correct
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return secret, nil
	})

	if err != nil {
		return 0, err
	}

	// Extract the claims (payload)
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// Extract the user ID from the claims
		if userIDStr, ok := claims[string(userIDKey)].((string)); ok {
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
