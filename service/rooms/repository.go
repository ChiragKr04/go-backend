package rooms

import "database/sql"

type RoomsRepository struct {
	db *sql.DB
}

// CreateRoom implements types.RoomRepository.
func (r *RoomsRepository) CreateRoom() error {
	return nil
}

func NewRepository(db *sql.DB) *RoomsRepository {
	return &RoomsRepository{
		db: db,
	}
}
