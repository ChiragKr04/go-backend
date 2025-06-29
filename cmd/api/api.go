package api

import (
	"ChiragKr04/go-backend/service/chat"
	"ChiragKr04/go-backend/service/rooms"
	"ChiragKr04/go-backend/service/user"
	"ChiragKr04/go-backend/service/webrtc"
	"ChiragKr04/go-backend/service/websocket"
	"database/sql"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type APIServerModel struct {
	address string
	db      *sql.DB
}

func APIServer(address string, db *sql.DB) *APIServerModel {
	return &APIServerModel{
		address: address,
		db:      db,
	}
}

func (s *APIServerModel) Run() error {
	// Run the server
	router := mux.NewRouter()
	// allow all CORS requests
	// Enhanced CORS middleware
	router.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Set CORS headers for all requests
			origin := r.Header.Get("Origin")
			if origin != "" {
				w.Header().Set("Access-Control-Allow-Origin", origin)
			} else {
				w.Header().Set("Access-Control-Allow-Origin", "*")
			}

			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.Header().Set("Access-Control-Max-Age", "86400") // 24 hours

			// Handle preflight OPTIONS request
			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}

			next.ServeHTTP(w, r)
		})
	})
	subRouter := router.PathPrefix("/api/v1").Subrouter()

	/// Repositories

	userRepo := user.NewRepository(s.db)
	roomsRepo := rooms.NewRepository(s.db)
	chatRepo := chat.NewRepository(s.db)
	webrtcRepo := webrtc.NewRepository(s.db)
	/// Handlers

	userHandler := user.NewHandler(userRepo)
	websocketHandler := websocket.NewHandler(userRepo, roomsRepo, chatRepo, webrtcRepo)
	roomsHandler := rooms.NewHandler(roomsRepo, userRepo)
	chatHandler := chat.NewHandler(chatRepo)

	/// Routes

	userHandler.UserRoutes(subRouter)
	roomsHandler.RoomRoutes(subRouter)
	websocketHandler.WebsocketRoutes(subRouter)
	chatHandler.ChatRoutes(subRouter)

	log.Println("Starting server on", s.address)

	return http.ListenAndServe(s.address, router)
}
