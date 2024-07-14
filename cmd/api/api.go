package api

import (
	"ChiragKr04/go-backend/service/user"
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

	log.Println("Starting server on", s.address)

	return http.ListenAndServe(s.address, router)
}
