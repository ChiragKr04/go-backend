package types

import "time"

type RoomRepository interface {
	CreateRoom(user *User, payload RoomCreateRequest) (int64, error)
	GetRoomById(id int64) (*Room, error)
	GetRoomByRoomId(roomId string) (*Room, error)
}

type WebSocketMessage struct {
	Type      string `json:"type"`
	UserID    int    `json:"userId"`
	Username  string `json:"username"`
	Data      any    `json:"data"`
	Timestamp string `json:"timestamp"`
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

// SocketEventType represents different types of socket events
type SocketEventType int

const (
	UserCountEvent       SocketEventType = iota
	UserJoinedEvent      SocketEventType = iota
	UserLeftEvent        SocketEventType = iota
	SendMessageEvent     SocketEventType = iota
	ReceiveMessageEvent  SocketEventType = iota
	JoinRoomEvent        SocketEventType = iota
	LeaveRoomEvent       SocketEventType = iota
	AuthenticateEvent    SocketEventType = iota
	ConnectEvent         SocketEventType = iota
	ConnectErrorEvent    SocketEventType = iota
	DisconnectEvent      SocketEventType = iota
	ReconnectFailedEvent SocketEventType = iota
	MessageEvent         SocketEventType = iota
	MessageReceivedEvent SocketEventType = iota
	MessageHistoryEvent  SocketEventType = iota
	SystemMessageEvent   SocketEventType = iota
	NewMessageEvent      SocketEventType = iota
)

// String method to get string representation of the enum
func (s SocketEventType) String() string {
	// UserCount = "USER_COUNT",
	// UserJoined = "USER_JOINED",
	// UserLeft = "USER_LEFT",
	// SendMessage = "SEND_MESSAGE",
	// ReceiveMessage = "RECEIVE_MESSAGE",
	// JoinRoom = "JOIN_ROOM",
	// LeaveRoom = "LEAVE_ROOM",
	// Authenticate = "AUTHENTICATE",
	// Connect = "CONNECT",
	// ConnectError = "CONNECT_ERROR",
	// Disconnect = "DISCONNECT",
	// ReconnectFailed = "RECONNECT_FAILED",
	// Message = "MESSAGE",
	// MessageReceived = "MESSAGE_RECEIVED",
	// MessageHistory = "MESSAGE_HISTORY",
	// SystemMessage = "SYSTEM_MESSAGE",
	// NewMessage = "NEW_MESSAGE",
	switch s {
	case UserCountEvent:
		return "USER_COUNT"
	case UserJoinedEvent:
		return "USER_JOINED"
	case UserLeftEvent:
		return "USER_LEFT"
	case SendMessageEvent:
		return "SEND_MESSAGE"
	case ReceiveMessageEvent:
		return "RECEIVE_MESSAGE"
	case JoinRoomEvent:
		return "JOIN_ROOM"
	case LeaveRoomEvent:
		return "LEAVE_ROOM"
	case AuthenticateEvent:
		return "AUTHENTICATE"
	case ConnectEvent:
		return "CONNECT"
	case ConnectErrorEvent:
		return "CONNECT_ERROR"
	case DisconnectEvent:
		return "DISCONNECT"
	case ReconnectFailedEvent:
		return "RECONNECT_FAILED"
	case MessageEvent:
		return "MESSAGE"
	case MessageReceivedEvent:
		return "MESSAGE_RECEIVED"
	case MessageHistoryEvent:
		return "MESSAGE_HISTORY"
	case SystemMessageEvent:
		return "SYSTEM_MESSAGE"
	case NewMessageEvent:
		return "NEW_MESSAGE"
	default:
		return "UNKNOWN"
	}
}
