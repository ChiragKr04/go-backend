package chat

import (
	"ChiragKr04/go-backend/types"
	"database/sql"
	"time"
)

type ChatRepository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *ChatRepository {
	return &ChatRepository{
		db: db,
	}
}

func (r *ChatRepository) RoomJoined(userId int, roomId string) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// check if same roodid with userid exists then dont insert the same user and roomid
	var count int
	err = tx.QueryRow("SELECT COUNT(*) FROM room_users WHERE user_id = ? AND room_id = ?", userId, roomId).Scan(&count)
	if err != nil {
		return err
	}
	if count > 0 {
		return nil
	}

	_, err = tx.Exec("INSERT INTO room_users (user_id, room_id) VALUES (?, ?)", userId, roomId)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}

func (r *ChatRepository) RoomLeft(userId int, roomId string) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.Exec("DELETE FROM room_users WHERE user_id = ? AND room_id = ?", userId, roomId)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}

func (r *ChatRepository) SaveChat(chat types.Chat) (types.Chat, error) {
	res, err := r.db.Exec(
		"INSERT INTO chats (userId, roomId, chat, chatType, createdAt) VALUES (?, ?, ?, ?, ?)",
		chat.UserID,
		chat.RoomID,
		chat.Chat,
		chat.ChatType,
		time.Now(),
	)
	if err != nil {
		return chat, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return chat, err
	}

	chat.ID = int(id)

	return chat, err
}

func (r *ChatRepository) GetChatsByRoomId(roomId string, limit int, offset int) ([]types.Chat, error) {
	rows, err := r.db.Query(
		"SELECT c.id, c.userId, c.roomId, c.chat, c.chatType, c.createdAt, u.username FROM chats c "+
			"LEFT JOIN users u ON c.userId = u.id "+
			"WHERE c.roomId = ? ORDER BY c.createdAt ASC LIMIT ? OFFSET ?",
		roomId, limit, offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var chats []types.Chat
	for rows.Next() {
		var chat types.Chat
		var createdAtStr string
		var username string
		err := rows.Scan(
			&chat.ID,
			&chat.UserID,
			&chat.RoomID,
			&chat.Chat,
			&chat.ChatType,
			&createdAtStr,
			&username,
		)
		if err != nil {
			return nil, err
		}

		createdAt, err := time.Parse(time.RFC3339, createdAtStr)
		if err != nil {
			return nil, err
		}
		chat.CreatedAt = createdAt
		chat.Username = username
		chats = append(chats, chat)
	}

	return chats, nil
}
