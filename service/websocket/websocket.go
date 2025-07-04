package websocket

import (
	"ChiragKr04/go-backend/service/rooms"
	"ChiragKr04/go-backend/service/user"
	"ChiragKr04/go-backend/service/webrtc"
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
	UserID   int
	UserName string
}

type WebhookHandler struct {
	userRepo   user.UserRepository
	webrtcRepo webrtc.WebrtcRepository
}

func NewWebhookHandler(userRepo user.UserRepository, webrtcRepo webrtc.WebrtcRepository) *WebhookHandler {
	return &WebhookHandler{
		userRepo:   userRepo,
		webrtcRepo: webrtcRepo,
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
		UserID:   userID,
		UserName: user.Username,
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
		var websocketMessage types.WebSocketMessage
		if err := json.Unmarshal(msg, &websocketMessage); err != nil {
			log.Printf("Error parsing message: %v", err)
			continue
		}

		// Set the user ID and room ID from the client
		websocketMessage.UserID = c.UserID
		websocketMessage.Username = c.UserName

		// Save the message to database if it's a text message
		if websocketMessage.Type == types.SendMessageEvent.String() {

			var chatMessage types.ChatMessage
			if err := json.Unmarshal(msg, &chatMessage); err != nil {
				log.Printf("Error parsing message: %v", err)
				continue
			}

			log.Printf("Received message: %v", chatMessage)
			chat := types.Chat{
				UserID:   c.UserID,
				RoomID:   c.RoomId,
				Chat:     chatMessage.Data.Chat,
				ChatType: "TEXT",
			}

			msgData, err := h.ChatRepo.SaveChat(chat)
			if err != nil {
				log.Printf("Error saving chat to database: %v", err)
			}
			msgData.Username = c.UserName

			marshaledData, err := json.Marshal(msgData)
			if err != nil {
				log.Printf("Error marshalling message: %v", err)
			}
			c.Hub.Broadcast <- marshaledData
			continue
		} else if websocketMessage.Type == types.JoinRoomEvent.String() {
			userData, err := h.ChatRepo.RoomJoined(c.UserID, c.RoomId)
			if err != nil {
				log.Printf("Error joining room: %v", err)
				continue
			}
			// log.Printf("User joined: %s", websocketMessage.Data.Username)
			log.Printf("Room users count: %d", len(userData))

			if len(userData) == 0 {
				log.Printf("No users in room")
				continue
			}

			msg, err := json.Marshal(types.WebSocketMessage{
				Type: types.UserCountEvent.String(),
				Data: types.RoomUserCountMessage{
					RoomUsersCount: len(userData),
					RoomUsersData:  userData,
				},
			})
			if err != nil {
				log.Printf("Error marshalling message: %v", err)
				continue
			}

			c.Hub.Broadcast <- msg
			continue
		} else if websocketMessage.Type == types.LeaveRoomEvent.String() {
			userData, err := h.ChatRepo.RoomLeft(c.UserID, c.RoomId)
			if err != nil {
				log.Printf("Error leaving room: %v", err)
				continue
			}
			// log.Printf("User left: %s", websocketMessage.Data.Username)
			log.Printf("Room users count: %d", len(userData))

			msg, err := json.Marshal(types.WebSocketMessage{
				Type: types.UserCountEvent.String(),
				Data: types.RoomUserCountMessage{
					RoomUsersCount: len(userData),
					RoomUsersData:  userData,
				},
			})
			if err != nil {
				log.Printf("Error marshalling message: %v", err)
				continue
			}
			c.Hub.Broadcast <- msg
			continue
		} else if websocketMessage.Type == types.NewOffer.String() {
			var offer types.WebRTCOfferMessage
			if err := json.Unmarshal(msg, &offer); err != nil {
				log.Printf("Error parsing message: %v", err)
				continue
			}
			// convert offer.Data.Offer.SDP to json string
			fullOfferJSON, err := json.Marshal(offer.Data.Offer)
			if err != nil {
				log.Printf("Error marshalling offer: %v", err)
				continue
			}
			offerData, err := h.WebrtcRepo.CreateOffer(types.Offer{
				Offer:               string(fullOfferJSON),
				Answer:              "",
				OffererUserID:       offer.UserID,
				AnswererUserID:      0,
				RoomID:              c.RoomId,
				OfferIceCandidates:  "",
				AnswerIceCandidates: "",
			})
			if err != nil {
				log.Printf("Error creating offer: %v", err)
				continue
			}
			msg, err := json.Marshal(types.WebSocketMessage{
				Type:     types.NewOfferAwaiting.String(),
				Data:     offerData,
				UserID:   offer.UserID,
				Username: offer.Username,
			})
			if err != nil {
				log.Printf("Error marshalling message: %v", err)
				continue
			}
			c.Hub.Broadcast <- msg
			continue
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
