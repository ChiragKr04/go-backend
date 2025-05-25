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

func (r *RoomsRepository) CreateRoom(user *types.User, payload types.RoomCreateRequest) (int64, error) {
	uuid := uuid.New()
	shortRoomID := generateShortRoomID()
	createdAt := time.Now()
	createdBy := user.ID
	tx, err := r.db.Begin()
	if err != nil {
		return -1, err
	}

	// Create a single invitation group for the room
	var invitationGroupID *int64 = nil
	if payload.Invitations != nil && len(payload.Invitations) > 0 {
		groupResult, err := tx.Exec("INSERT INTO invitation_groups (room_id, created_at) VALUES (?, ?)", uuid, createdAt)
		if err != nil {
			_ = tx.Rollback()
			return -1, err
		}
		invitationGroupIDValue, err := groupResult.LastInsertId()
		invitationGroupID = &invitationGroupIDValue
		if err != nil {
			_ = tx.Rollback()
			return -1, err
		}
		// Add all users to the invitation group
		for _, userID := range payload.Invitations {
			_, err = tx.Exec("INSERT INTO invitation_users (invitation_group_id, user_id) VALUES (?, ?)", invitationGroupID, userID)
			if err != nil {
				_ = tx.Rollback()
				return -1, err
			}
		}
	}

	result, err := tx.Exec("INSERT INTO rooms (room_id, short_room_id, created_at, created_by, is_private, invitations) VALUES (?, ?, ?, ?, ?, ?)", uuid, shortRoomID, createdAt, createdBy, payload.IsPrivate, invitationGroupID)
	if err != nil {
		_ = tx.Rollback()
		return -1, err
	}

	if err = tx.Commit(); err != nil {
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

	// Fetch invitation users for this room
	err = r.fetchInvitationUsers(room)
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

	// Fetch invitation users for this room
	if room.ID != 0 {
		err = r.fetchInvitationUsers(room)
		if err != nil {
			return nil, err
		}
	}

	if room.ID == 0 {
		return nil, fmt.Errorf("room not found")
	}
	return room, nil
}

func scanRowsIntoRoom(rows *sql.Rows) (*types.Room, error) {
	room := &types.Room{}
	for rows.Next() {
		// Need to scan all 7 columns from the rooms table
		// The 7th column is 'invitations' which is stored as INT in the database
		// but we'll handle it separately with fetchInvitationUsers
		var dbInvitation sql.NullInt64 // Use NullInt64 to handle NULL values
		err := rows.Scan(&room.ID, &room.RoomId, &room.CreatedBy, &room.CreatedAt, &room.IsPrivate, &room.ShortRoomId, &dbInvitation)
		if err != nil {
			return nil, err
		}
		// We don't use the dbInvitation value as we'll populate Invitations using fetchInvitationUsers
	}
	return room, nil
}

// fetchInvitationUsers retrieves all users invited to a room and populates the room's Invitations field
func (r *RoomsRepository) fetchInvitationUsers(room *types.Room) error {
	// Query to get all users associated with invitation groups for this room
	query := `
		SELECT iu.user_id 
		FROM invitation_groups ig
		JOIN invitation_users iu ON ig.id = iu.invitation_group_id
		WHERE ig.room_id = ?
	`

	rows, err := r.db.Query(query, room.RoomId)
	if err != nil {
		return err
	}
	defer rows.Close()

	// Initialize the invitations slice
	room.Invitations = []*types.User{}

	// Collect all user IDs
	userIDs := []int{}
	for rows.Next() {
		var userID int
		err := rows.Scan(&userID)
		if err != nil {
			return err
		}
		userIDs = append(userIDs, userID)
	}

	// Fetch complete user data for each user ID
	for _, userID := range userIDs {
		// Query to get user data
		userQuery := `
			SELECT id, first_name, last_name, email, password, createdAt
			FROM users
			WHERE id = ?
		`
		userRow := r.db.QueryRow(userQuery, userID)

		user := &types.User{}
		err := userRow.Scan(
			&user.ID,
			&user.FirstName,
			&user.LastName,
			&user.Email,
			&user.Password,
			&user.CreatedAt,
		)
		if err != nil {
			return err
		}

		// For security, clear the password before adding to response
		user.Password = ""

		room.Invitations = append(room.Invitations, user)
	}

	return nil
}

func NewRepository(db *sql.DB) *RoomsRepository {
	return &RoomsRepository{
		db: db,
	}
}
