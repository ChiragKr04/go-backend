package api

import (
	"ChiragKr04/go-backend/service/rooms"
	"ChiragKr04/go-backend/service/user"
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
	subRouter := router.PathPrefix("/api/v1").Subrouter()

	/// Repositories

	userRepo := user.NewRepository(s.db)
	roomsRepo := rooms.NewRepository(s.db)
	
	/// Handlers

	userHandler := user.NewHandler(userRepo)
	websocketHandler := websocket.NewHandler(userRepo, roomsRepo)
	roomsHandler := rooms.NewHandler(roomsRepo, userRepo)
	
	/// Routes
	
	userHandler.UserRoutes(subRouter)
	roomsHandler.RoomRoutes(subRouter)
	websocketHandler.WebsocketRoutes(subRouter)

	
	log.Println("Starting server on", s.address)

	return http.ListenAndServe(s.address, router)
}
