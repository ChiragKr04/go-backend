package rooms

import (
	"ChiragKr04/go-backend/types"
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type RoomsRepository struct {
	db *sql.DB
}

func (r *RoomsRepository) CreateRoom(user *types.User) (int64, error) {
	uuid := uuid.New()
	createdAt := time.Now()
	createdBy := user.ID
	result, err := r.db.Exec("INSERT INTO rooms (room_id, created_at, created_by) VALUES (?, ?, ?)", uuid, createdAt, createdBy)
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
	return room, nil
}

func scanRowsIntoRoom(rows *sql.Rows) (*types.Room, error) {
	room := &types.Room{}
	for rows.Next() {
		err := rows.Scan(&room.ID, &room.RoomId, &room.CreatedBy, &room.CreatedAt)
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
