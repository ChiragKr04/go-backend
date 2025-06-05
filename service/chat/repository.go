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
