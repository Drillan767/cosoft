package storage

import (
	"database/sql"
	"errors"
	"time"
)

type Store struct {
	db *sql.DB
}

func NewStore(dbPath string) (*Store, error) {
	db, err := sql.Open("sqlite3", dbPath)

	if err != nil {
		return nil, err
	}

	return &Store{db: db}, nil
}

func (s *Store) SetupDatabase(dbPath string) error {
	db, err := sql.Open("sqlite3", dbPath)

	if err != nil {
		return err
	}

	defer db.Close()

	query := `
		CREATE TABLE IF NOT EXISTS users (
			id VARCHAR(40) NOT NULL,
			first_name VARCHAR(50) NOT NULL,
			last_name VARCHAR(50) NOT NULL,
			email VARCHAR(50) NOT NULL,
			credits SMALLINT NOT NULL DEFAULT 0,
			jwt VARCHAR(50) NOT NULL,
			expires_at DATE NOT NULL,
			slack_user_id VARCHAR(50),
			created_at DATE NOT NULL
		);

		CREATE TABLE IF NOT EXISTS rooms (
			id VARCHAR(40) NOT NULL,
			name VARCHAR(50) NOT NULL,
			nb_users TINYINT NOT NULL DEFAULT 0,
			price SMALLINT NOT NULL DEFAULT 0,
			created_at DATE NOT NULL
		);
	`

	_, err = s.db.Exec(query)

	if err != nil {
		return err
	}

	return nil
}

func (s *Store) GetUserByEmail(email string) (*User, error) {
	// Check if JWT is still valid
	query := `SELECT * FROM users WHERE email = ? AND expires_at > ?`
	var user User
	err := s.db.QueryRow(query, email, time.Now()).Scan(&user)
	if err == sql.ErrNoRows {
		return nil, errors.New("no valid session found")
	}
	return &user, err
}

func (s *Store) CreateUser(user *User) error {
	query := `
        INSERT INTO users (id, email, first_name, last_name, credits, jwt, expires_at, slack_user_id, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT (email) DO UPDATE SET
			id = excluded.id,
			email = excluded.email,
			first_name = excluded.first_name,
			last_name = excluded.last_name,
			credits = excluded.credits,
			jwt = excluded.jwt,
			expires_at = excluded.expires_at,
			slack_user_id = excluded.slack_user_id,
			created_at = excluded.created_at
		`

	_, err := s.db.Exec(
		query,
		user.Id,
		user.Email,
		user.FirstName,
		user.LastName,
		user.Credits,
		user.Jwt,
		user.ExpiresAt,
		user.SlackUserID,
		time.Now(),
	)

	return err
}

// Update credits

// List reservations (parameter: paste / future)
// Store reservation
// Remove (cancel) reservation
