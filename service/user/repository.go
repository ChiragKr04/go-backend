package user

import (
	"ChiragKr04/go-backend/types"
	"ChiragKr04/go-backend/utils"
	"database/sql"
	"fmt"
	"time"
)

type UserRepository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r *UserRepository) GetUserByEmail(email string) (*types.User, error) {
	rows, err := r.db.Query(("SELECT * FROM users WHERE email = ?"), email)
	if err != nil {
		return nil, err
	}
	user := &types.User{}
	for rows.Next() {
		user, err = scanRowsIntoUser(rows)
		if err != nil {
			return nil, err
		}
	}
	if user.ID == 0 {
		return nil, fmt.Errorf("user not found")
	}
	return user, nil
}

func (r *UserRepository) CreateUser(user types.User) error {
	_, err := r.db.Exec(
		"INSERT INTO users (first_name, last_name, email, password, username) VALUES (?, ?, ?, ?, ?)",
		user.FirstName,
		user.LastName,
		user.Email,
		user.Password,
		user.Username,
	)
	if err != nil {
		return err
	}
	return nil
}

func (r *UserRepository) GetUserByID(id int) (*types.User, error) {
	rows, err := r.db.Query(("SELECT * FROM users WHERE id = ?"), id)
	if err != nil {
		return nil, err
	}
	user := &types.User{}
	for rows.Next() {
		user, err = scanRowsIntoUser(rows)
		if err != nil {
			return nil, err
		}
	}
	if user.ID == 0 {
		return nil, fmt.Errorf("user not found")
	}
	return user, nil
}

func (r *UserRepository) UpdateUser(user types.User) (*types.User, error) {
	_, err := r.db.Exec(
		"UPDATE users SET first_name = ?, last_name = ?, email = ? WHERE id = ?",
		user.FirstName,
		user.LastName,
		user.Email,
		user.ID,
	)
	if err != nil {
		return nil, err
	}
	return r.GetUserByID(user.ID)
}

func scanRowsIntoUser(rows *sql.Rows) (*types.User, error) {
	user := &types.User{}
	var createdAtStr string
	if err := rows.Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.Password,
		&createdAtStr,
		&user.Username,
	); err != nil {
		return nil, err
	}

	createdAt, err := time.Parse(time.RFC3339, createdAtStr)
	if err != nil {
		return nil, err
	}
	user.CreatedAt = createdAt
	return user, nil
}

func (r *UserRepository) SearchUser(query string) ([]types.User, error) {
	rows, err := r.db.Query(("SELECT * FROM users WHERE username LIKE ? OR email LIKE ?"), "%"+query+"%", "%"+query+"%")
	if err != nil {
		return nil, err
	}
	users := []types.User{}
	for rows.Next() {
		user, err := scanRowsIntoUser(rows)
		if err != nil {
			return nil, err
		}
		tempUser := utils.ReturnUserWithoutPassword(*user)
		users = append(users, tempUser)
	}
	return users, nil
}
