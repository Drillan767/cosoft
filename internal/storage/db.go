package storage

import (
	"cosoft-cli/internal/api"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3"
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

	query := `
		CREATE TABLE IF NOT EXISTS users (
			id VARCHAR(40) PRIMARY KEY NOT NULL,
			first_name VARCHAR(50) NOT NULL,
			last_name VARCHAR(50) NOT NULL,
			email VARCHAR(50) UNIQUE NOT NULL,
			credits SMALLINT NOT NULL DEFAULT 0,
			jwt VARCHAR(50) NOT NULL,
			expires_at DATE NOT NULL,
			slack_user_id VARCHAR(50),
			created_at DATE NOT NULL
		);

		CREATE TABLE IF NOT EXISTS rooms (
			id VARCHAR(40) PRIMARY KEY NOT NULL,
			name VARCHAR(50) NOT NULL,
			nb_users TINYINT NOT NULL DEFAULT 0,
			price SMALLINT NOT NULL DEFAULT 0,
			created_at DATE NOT NULL
		);
	`

	_, err := s.db.Exec(query)

	if err != nil {
		return err
	}

	return nil
}

func (s *Store) HasActiveToken() (bool, error) {
	// Check if JWT is still valid
	query := `SELECT 1 FROM users WHERE expires_at > ?`

	var result int
	err := s.db.QueryRow(query, time.Now()).Scan(&result)

	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

func (s *Store) CreateUser(user *api.UserResponse, expiresAt time.Time) error {
	var nbUsers int

	rows, err := s.db.Query("SELECT COUNT(*) FROM users;")

	if err != nil {
		return err
	}

	defer rows.Close()

	for rows.Next() {
		if err := rows.Scan(&nbUsers); err != nil {
			return err
		}
	}

	if nbUsers > 0 {
		return fmt.Errorf("A user already exists, aborting.")
	}

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

	_, err = s.db.Exec(
		query,
		user.Id,
		user.Email,
		user.FirstName,
		user.LastName,
		user.Credits*100,
		user.JwtToken,
		expiresAt,
		nil,
		time.Now(),
	)

	return err
}

func (s *Store) GetUserData() (*User, error) {

	var user User

	query := `
		SELECT id, first_name, last_name, email, jwt, credits, expires_at, slack_user_id, created_at FROM users LIMIT 1;
	`

	err := s.db.QueryRow(query).Scan(
		&user.Id,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.Jwt,
		&user.Credits,
		&user.ExpiresAt,
		&user.SlackUserID,
		&user.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *Store) LogoutUser() error {
	query := `
		UPDATE TABLE users SET jwt = "", expires_at = ? WHERE id = (SELECT id FROM users LIMIT 1)
	`

	_, err := s.db.Exec(query, time.Now())

	return err
}

// Update credits

// List reservations (parameter: paste / future)
// Store reservation
// Remove (cancel) reservation
