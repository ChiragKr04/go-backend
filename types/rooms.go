package types

import "time"

type RoomRepository interface {
	CreateRoom(user *User) (int64, error)
	GetRoomById(id int64) (*Room, error)
	GetRoomByRoomId(roomId string) (*Room, error)
}

type Room struct {
	ID        int       `json:"id"`
	RoomId    string    `json:"room_id"`
	CreatedAt time.Time `json:"created_at"`
	CreatedBy int       `json:"created_by"`
}
