package websocket

import (
	"ChiragKr04/go-backend/types"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

type HubFile struct {
	HubType *types.Hub
}

func (h *HubFile) Run() {
	for {
		select {
		case client := <-h.HubType.Register:
			h.HubType.Clients[client] = true
		case client := <-h.HubType.Unregister:
			if _, ok := h.HubType.Clients[client]; ok {
				delete(h.HubType.Clients, client)
				close(client.Send)
			}
		case message := <-h.HubType.Broadcast:
			for client := range h.HubType.Clients {
				select {
				case client.Send <- message:
				default:
					close(client.Send)
					delete(h.HubType.Clients, client)
				}
			}
		}
	}
}

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

var (
	newline = []byte{'\n'}
	// space   = []byte{' '}
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// Allow all connections by returning true
		return true
	},
}
