package rooms

import (
	"ChiragKr04/go-backend/types"
	"database/sql"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/google/uuid"
)

type RoomsRepository struct {
	db *sql.DB
}

// generateShortRoomID creates a short room ID like "ZAQ166"
func generateShortRoomID() string {
	const chars = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	const numbers = "0123456789"
	
	rand.Seed(time.Now().UnixNano())
	
	result := strings.Builder{}
	
	// Add 3 random letters
	for i := 0; i < 3; i++ {
		result.WriteByte(chars[rand.Intn(len(chars))])
	}
	
	// Add 3 random numbers
	for i := 0; i < 3; i++ {
		result.WriteByte(numbers[rand.Intn(len(numbers))])
	}
	
	return result.String()
}

func (r *RoomsRepository) CreateRoom(user *types.User, isPrivate bool) (int64, error) {
	uuid := uuid.New()
	shortRoomID := generateShortRoomID()
	createdAt := time.Now()
	createdBy := user.ID
	result, err := r.db.Exec("INSERT INTO rooms (room_id, short_room_id, created_at, created_by, is_private) VALUES (?, ?, ?, ?, ?)", uuid, shortRoomID, createdAt, createdBy, isPrivate)
	if err != nil {
		return -1, err
	}

	roomData, err := result.LastInsertId()
	if err != nil {
		return -1, err
	}
	return roomData, nil
}

func (r *RoomsRepository) GetRoomById(id int64) (*types.Room, error) {
	res, err := r.db.Query("SELECT * FROM rooms WHERE id = ?", id)
	if err != nil {
		return nil, err
	}

	room, err := scanRowsIntoRoom(res)
	if err != nil {
		return nil, err
	}
	return room, nil
}

func (r *RoomsRepository) GetRoomByRoomId(roomId string) (*types.Room, error) {
	res, err := r.db.Query("SELECT * FROM rooms WHERE room_id = ?", roomId)
	if err != nil {
		return nil, err
	}
	room, err := scanRowsIntoRoom(res)
	if err != nil {
		return nil, err
	}
	if room.ID == 0 {
		return nil, fmt.Errorf("room not found")
	}
	return room, nil
}

func scanRowsIntoRoom(rows *sql.Rows) (*types.Room, error) {
	room := &types.Room{}
	for rows.Next() {
		err := rows.Scan(&room.ID, &room.RoomId, &room.CreatedBy, &room.CreatedAt, &room.IsPrivate, &room.ShortRoomId)
		if err != nil {
			return nil, err
		}
	} 
	return room, nil
}

func NewRepository(db *sql.DB) *RoomsRepository {
	return &RoomsRepository{
		db: db,
	}
}
