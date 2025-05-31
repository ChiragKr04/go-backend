package websocket

import (
	"ChiragKr04/go-backend/service/rooms"
	"ChiragKr04/go-backend/service/user"
	"ChiragKr04/go-backend/types"
	"ChiragKr04/go-backend/utils"

	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

type LocalClient struct {
	*types.Client
	UserID int
}

type WebhookHandler struct {
	userRepo user.UserRepository
}

func NewWebhookHandler(userRepo user.UserRepository) *WebhookHandler {
	return &WebhookHandler{
		userRepo: userRepo,
	}
}

// serveWs handles websocket requests from the peer.
func (h *Handler) serveWs(w http.ResponseWriter, r *http.Request) {
	// take token in query
	token := r.URL.Query().Get("token")
	if token == "" {
		http.Error(w, "Token is required", http.StatusBadRequest)
		return
	}
	userID, err := utils.GetUserIDFromToken(w, r)
	if err != nil || userID == 0 {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}
	user, err := h.UserRepo.GetUserByID(userID)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}
	print(user.Email)
	room := GetValidRoom(r, w, h.RoomRepo)
	if room == nil {
		http.Error(w, "Room ID not found", http.StatusBadRequest)
		return
	}

	// Get the hub for the room
	hub := h.HubManager.GetHub(room.RoomId)

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}
	//
	client := &LocalClient{
		Client: &types.Client{
			Hub:     hub.HubType,
			Send:    make(chan []byte, 256),
			Request: r,
			RoomId:  room.RoomId,
		},
		UserID: userID,
	}
	client.Conn = conn

	// Register the client with the hub
	hub.HubType.Register <- client.Client

	// Start the pumps
	go h.WritePump(client)
	go h.ReadPump(client)
}

func (h *Handler) ReadPump(c *LocalClient) {
	defer func() {
		c.Hub.Unregister <- c.Client
		c.Conn.Close()
	}()
	// c.Conn.SetReadLimit(maxMessageSize)
	// c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	// c.Conn.SetPongHandler(func(string) error { c.Conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, msg, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}

		// Parse the incoming message
		var chatMessage types.ChatMessage
		if err := json.Unmarshal(msg, &chatMessage); err != nil {
			log.Printf("Error parsing message: %v", err)
			continue
		}

		// Set the user ID and room ID from the client
		chatMessage.UserID = c.UserID

		// Save the message to database if it's a text message
		if chatMessage.Type == "send_message" {
			chat := types.Chat{
				UserID:   c.UserID,
				RoomID:   c.RoomId,
				Chat:     chatMessage.Payload.Chat,
				ChatType: "TEXT",
			}

			if err := h.ChatRepo.SaveChat(chat); err != nil {
				log.Printf("Error saving chat to database: %v", err)
			}
		}

		// Broadcast the message to all clients in the room
		c.Hub.Broadcast <- msg
	}
}

func (h *Handler) WritePump(c *LocalClient) {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.Send:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			n := len(c.Send)
			for i := 0; i < n; i++ {
				w.Write(newline)
				w.Write(<-c.Send)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func GetValidRoom(r *http.Request, w http.ResponseWriter, roomRepo *rooms.RoomsRepository) *types.Room {
	vars := mux.Vars(r)
	roomId := vars["roomId"]

	log.Printf("Room ID: %s", roomId)

	room, err := roomRepo.GetRoomByRoomId(roomId)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return nil
	}

	return room
}
