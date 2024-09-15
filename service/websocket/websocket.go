package websocket

import (
	"ChiragKr04/go-backend/service/user"
	"ChiragKr04/go-backend/types"
	// "bytes"
	"fmt"

	// "encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

type LocalClient struct {
	*types.Client
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
func (h *Handler) serveWs(hub *HubFile, w http.ResponseWriter, r *http.Request) {

	client := &LocalClient{
		Client: &types.Client{
			Hub:            hub.HubType,
			Send:           make(chan []byte, 256),
			ResponseWriter: w,
			Request:        r,
		},
	}

	user := GetUserData(client.Client, h.UserRepo)
	if user == nil {
		http.Error(w, "User ID not found", http.StatusBadRequest)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	client.Conn = conn
	client.Hub.Register <- client.Client
	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go h.WritePump(client)
	go h.ReadPump(client)
}

func (h *Handler) ReadPump(c *LocalClient) {
	defer func() {
		c.Hub.Unregister <- c.Client
		c.Conn.Close()
	}()
	c.Conn.SetReadLimit(maxMessageSize)
	c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	c.Conn.SetPongHandler(func(string) error { c.Conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, msg, err := c.Conn.ReadMessage()
		fmt.Println(string(msg))
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		// user := GetUserData(c.Client, h.UserRepo)
		// userJson, err := json.Marshal(
		// 	map[string]interface{}{
		// 		"message": string(msg),
		// 		"user":    user,
		// 	},
		// )
		// if err != nil {
		// 	log.Println(err)
		// }
		// message := bytes.TrimSpace(bytes.Replace(msg, nil, nil, -1))
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
				// The hub closed the channel.
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

func GetUserData(c *types.Client, userRepo *user.UserRepository) *types.User {
	vars := mux.Vars(c.Request)
	userIDStr := vars["userId"]

	log.Printf("User ID: %s", userIDStr)

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		http.Error(c.ResponseWriter, "Invalid user ID", http.StatusBadRequest)
		return nil
	}

	if userID == 0 {
		log.Println("User ID not found")
		http.Error(c.ResponseWriter, "User ID not found", http.StatusBadRequest)
	}

	user, err := userRepo.GetUserByID(userID)
	if err != nil {
		http.Error(c.ResponseWriter, "User not found", http.StatusNotFound)
		return nil
	}

	return user
}
