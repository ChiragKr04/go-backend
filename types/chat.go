package types

import (
	"time"
)

type ChatRepository interface {
	SaveChat(chat Chat) (Chat, error)
	GetChatsByRoomId(roomId string, limit int, offset int) ([]Chat, error)
	RoomJoined(userId int, roomId string) ([]RoomUserData, error)
	RoomLeft(userId int, roomId string) ([]RoomUserData, error)
	GetRoomUsersCount(roomId string) ([]RoomUserData, error)
}

type Chat struct {
	ID        int       `json:"id"`
	UserID    int       `json:"userId"`
	RoomID    string    `json:"roomId"`
	Chat      string    `json:"chat"`
	ChatType  string    `json:"chatType"`
	CreatedAt time.Time `json:"createdAt"`
	Username  string    `json:"username,omitempty"`
}

type ChatPayload struct {
	Chat      string `json:"chat"`
	Timestamp string `json:"timestamp"`
}

type ChatMessage struct {
	Type      string      `json:"type"`
	Data      ChatPayload `json:"data"`
	Timestamp string      `json:"timestamp"`
	UserID    int         `json:"userId"`
	Username  string      `json:"username"`
}

type RoomUserCountMessage struct {
	Type           string         `json:"type"`
	RoomUsersCount int            `json:"roomUsersCount"`
	RoomUsersData  []RoomUserData `json:"roomUsersData"`
}

type RoomUserData struct {
	UserID   int    `json:"userId"`
	Username string `json:"username"`
}
