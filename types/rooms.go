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
	UserCountEvent                    SocketEventType = iota
	UserJoinedEvent                   SocketEventType = iota
	UserLeftEvent                     SocketEventType = iota
	SendMessageEvent                  SocketEventType = iota
	ReceiveMessageEvent               SocketEventType = iota
	JoinRoomEvent                     SocketEventType = iota
	LeaveRoomEvent                    SocketEventType = iota
	AuthenticateEvent                 SocketEventType = iota
	ConnectEvent                      SocketEventType = iota
	ConnectErrorEvent                 SocketEventType = iota
	DisconnectEvent                   SocketEventType = iota
	ReconnectFailedEvent              SocketEventType = iota
	MessageEvent                      SocketEventType = iota
	MessageReceivedEvent              SocketEventType = iota
	MessageHistoryEvent               SocketEventType = iota
	SystemMessageEvent                SocketEventType = iota
	NewMessageEvent                   SocketEventType = iota
	SendIceCandidateToSignalingServer SocketEventType = iota
	NewOffer                          SocketEventType = iota
	NewOfferAwaiting                  SocketEventType = iota
	NewAnswer                         SocketEventType = iota
	AnswerResponse                    SocketEventType = iota
	ReceivedIceCandidateFromServer    SocketEventType = iota
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
	case SendIceCandidateToSignalingServer:
		return "SEND_ICE_CANDIDATE_TO_SIGNALING_SERVER"
	case NewOffer:
		return "NEW_OFFER"
	case NewOfferAwaiting:
		return "NEW_OFFER_AWAITING"
	case NewAnswer:
		return "NEW_ANSWER"
	case AnswerResponse:
		return "ANSWER_RESPONSE"
	case ReceivedIceCandidateFromServer:
		return "RECEIVED_ICE_CANDIDATE_FROM_SERVER"
	default:
		return "UNKNOWN"
	}
}

// type Offer struct {
// 	// offererUserName: userName,
// 	//           offer: newOffer,
// 	//           offerIceCandidates: [],
// 	//           answererUserName: null,
// 	//           answer: null,
// 	//           answererIceCandidates: []
// 	OffererUserName       string    `json:"offerer_user_name"`
// 	Offer                 OfferData `json:"offer"`
// 	OfferIceCandidates    []string  `json:"offer_ice_candidates"`
// 	AnswererUserName      *string   `json:"answerer_user_name"`
// 	Answer                *string   `json:"answer"`
// 	AnswererIceCandidates []string  `json:"answerer_ice_candidates"`
// }

type WebRTCOfferMessage struct {
	Type     string `json:"type"`
	UserID   int    `json:"userId"`
	Username string `json:"username"`
	Data     struct {
		Offer struct {
			SDP  string `json:"sdp"`
			Type string `json:"type,omitempty"`
		} `json:"offer"`
	} `json:"data"`
}

type Offer struct {
	Offer               string `json:"offer"`
	Answer              string `json:"answer"`
	OffererUserID       int    `json:"offerer_user_id"`
	AnswererUserID      int    `json:"answerer_user_id"`
	RoomID              string `json:"room_id"`
	OfferIceCandidates  string `json:"offer_ice_candidates"`
	AnswerIceCandidates string `json:"answer_ice_candidates"`
}

type OfferData struct {
	Type string `json:"type"`
	Sdp  any    `json:"sdp"`
}
