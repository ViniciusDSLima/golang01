package user

import (
	"database/sql"
	"fmt"
	"github.com/ViniciusDSLima/golang01/types"
	_ "net"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

func (s *Store) CreateUser(user types.User) error {
	_, err := s.db.Exec("INSERT INTO users (firstName, lastName, email, password) VALUES ($1,$2,$3,$4)",
		user.FirstName, user.LastName, user.Email, user.Password)

	if err != nil {
		return err
	}

	return nil
}

func (s *Store) GetUsersByEmail(email string) (*types.User, error) {

	rows, err := s.db.Query(`SELECT * FROM users WHERE "email" = $1`, email)

	if err != nil {
		return nil, err
	}

	u := new(types.User)

	for rows.Next() {
		u, err = scanRowIntoUser(rows)

		if err != nil {
			return nil, err
		}
	}

	if u.Id == "" {
		return nil, fmt.Errorf("user not found")
	}

	return u, nil
}

func (s *Store) GetUserById(id string) (*types.User, error) {
	rows, err := s.db.Query("SELECT * FROM users WHERE id = $1", id)

	if err != nil {
		return nil, err
	}

	u := new(types.User)

	for rows.Next() {
		_, err := scanRowIntoUser(rows)
		if err != nil {
			return nil, err
		}
	}

	if u.Id == "" {
		return nil, fmt.Errorf("user not found")
	}

	return u, nil

}
func scanRowIntoUser(rows *sql.Rows) (*types.User, error) {
	user := new(types.User)

	err := rows.Scan(
		&user.Id,
		&user.Email,
		&user.FirstName,
		&user.LastName,
		&user.CreatedAt,
		&user.Password,
	)

	if err != nil {
		return nil, err
	}

	return user, nil
}
