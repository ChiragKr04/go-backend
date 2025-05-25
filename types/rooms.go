package types

import "time"

type RoomRepository interface {
	CreateRoom(user *User, payload RoomCreateRequest) (int64, error)
	GetRoomById(id int64) (*Room, error)
	GetRoomByRoomId(roomId string) (*Room, error)
}

type Room struct {
	ID              int       `json:"id"`
	RoomId          string    `json:"room_id"`
	ShortRoomId     string    `json:"short_room_id"`
	CreatedAt       time.Time `json:"created_at"`
	CreatedBy       int       `json:"created_by"`
	IsPrivate       bool      `json:"is_private"`
	Invitations     []*User   `json:"invitations"`
	RoomName        string    `json:"room_name"`
	RoomDescription string    `json:"room_description"`
}

type RoomCreateRequest struct {
	RoomName        string `json:"room_name" validate:"required,max=255"`
	RoomDescription string `json:"room_description" validate:"required,max=1000"`
	IsPrivate       bool   `json:"is_private"`
	Invitations     []int  `json:"invitations" validate:"omitempty,min=1"`
}
