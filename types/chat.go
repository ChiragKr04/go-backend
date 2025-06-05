package types

import (
	"time"
)

type ChatRepository interface {
	SaveChat(chat Chat) (Chat, error)
	GetChatsByRoomId(roomId string, limit int, offset int) ([]Chat, error)
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
	UserID    int    `json:"userId"`
	Username  string `json:"username"`
}

type ChatMessage struct {
	Type      string      `json:"type"`
	Chat_Data ChatPayload `json:"chat_data"`
	Timestamp string      `json:"timestamp"`
	UserID    int         `json:"userId"`
	Username  string      `json:"username"`
}
