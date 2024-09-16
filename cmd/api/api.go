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

	userRepo := user.NewRepository(s.db)
	userHandler := user.NewHandler(userRepo)
	userHandler.UserRoutes(subRouter)

	websocketHandler := websocket.NewHandler(userRepo)
	websocketHandler.WebsocketRoutes(subRouter)

	roomsRepo := rooms.NewRepository(s.db)
	roomsHandler := rooms.NewHandler(roomsRepo)
	roomsHandler.RoomRoutes(subRouter)

	log.Println("Starting server on", s.address)

	return http.ListenAndServe(s.address, router)
}
