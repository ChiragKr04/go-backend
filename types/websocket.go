package types

import (
	"net/http"

	"github.com/gorilla/websocket"
)

type Client struct {
	Hub *Hub

	// The websocket connection.
	Conn *websocket.Conn

	// Buffered channel of outbound messages.
	Send chan []byte

	Request *http.Request

	RoomId string
}

type Hub struct {
	// Registered Clients.
	Clients map[*Client]bool

	// Inbound messages from the clients.
	Broadcast chan []byte

	// Register requests from the clients.
	Register chan *Client

	// Unregister requests from clients.
	Unregister chan *Client
}
