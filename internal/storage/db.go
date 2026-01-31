package storage

import (
	"cosoft-cli/internal/api"
	"cosoft-cli/shared/models"
	"database/sql"
	"errors"
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type Store struct {
	db *sql.DB
}

type Cookies struct {
	WAuth        string `db:"w_auth"`
	WAuthRefresh string `db:"w_auth_refresh"`
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
			credits REAL NOT NULL DEFAULT 0,
			w_auth TEXT NOT NULL,
			w_auth_refresh TEXT NOT NULL,
			slack_user_id VARCHAR(50),
			created_at DATE NOT NULL
		);

		CREATE TABLE IF NOT EXISTS rooms (
			id VARCHAR(40) PRIMARY KEY NOT NULL,
			name VARCHAR(50) NOT NULL,
			nb_users TINYINT NOT NULL DEFAULT 0,
			price REAL NOT NULL DEFAULT 0,
			created_at DATE NOT NULL
		);
	`

	_, err := s.db.Exec(query)

	if err != nil {
		return err
	}

	return nil
}

func (s *Store) HasActiveToken() (*Cookies, error) {
	query := `SELECT w_auth, w_auth_refresh FROM users LIMIT 1`

	var result Cookies
	err := s.db.QueryRow(query).Scan(
		&result.WAuth,
		&result.WAuthRefresh,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &result, nil
}

func (s *Store) SetUser(user *api.UserResponse, wAuth, wAuthRefresh string) error {
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

	if nbUsers == 0 {
		query := `
		        INSERT INTO users (id, email, first_name, last_name, credits, w_auth, w_auth_refresh, slack_user_id, created_at)
				VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
				ON CONFLICT (email) DO UPDATE SET
					id = excluded.id,
					email = excluded.email,
					first_name = excluded.first_name,
					last_name = excluded.last_name,
					credits = excluded.credits,
					w_auth = excluded.w_auth,
					w_auth_refresh = excluded.w_auth_refresh,
					slack_user_id = excluded.slack_user_id,
					created_at = excluded.created_at
				`

		_, err = s.db.Exec(
			query,
			user.Id,
			user.Email,
			user.FirstName,
			user.LastName,
			user.Credits,
			wAuth,
			wAuthRefresh,
			nil,
			time.Now(),
		)

		return err
	} else if nbUsers == 1 {
		query := `UPDATE users SET w_auth = ?, w_auth_refresh = ? WHERE id = ?`
		_, err := s.db.Exec(query, wAuth, wAuthRefresh, user.Id)

		return err
	}

	return errors.New("too many users")
}

func (s *Store) GetUserData() (*User, error) {

	var user User

	query := `
		SELECT id, first_name, last_name, email, w_auth, w_auth_refresh, credits, slack_user_id, created_at FROM users LIMIT 1;
	`

	err := s.db.QueryRow(query).Scan(
		&user.Id,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.WAuth,
		&user.WAuthRefresh,
		&user.Credits,
		&user.SlackUserID,
		&user.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *Store) LogoutUser() error {
	query := `DELETE FROM users`

	_, err := s.db.Exec(query)

	return err
}

func (s *Store) UpdateCredits() (*float64, error) {
	type UserCookies struct {
		Id      string  `db:"id"`
		Credits float64 `db:"credits"`
		Auth    string  `db:"w_auth"`
		Refresh string  `db:"w_auth_refresh"`
	}

	uc := UserCookies{}

	query := `
		SELECT id, credits, w_auth, w_auth_refresh FROM users LIMIT 1
	`

	err := s.db.QueryRow(query).Scan(
		&uc.Id,
		&uc.Credits,
		&uc.Auth,
		&uc.Refresh,
	)

	if err != nil {
		return nil, err
	}

	clientApi := api.NewApi()
	newCredits, err := clientApi.GetCredits(uc.Auth, uc.Refresh)

	if err != nil {
		return nil, err
	}

	if newCredits == uc.Credits {
		return nil, nil
	}

	query = `UPDATE users SET credits = ? WHERE id = ?`

	_, err = s.db.Exec(query, newCredits, uc.Id)

	return &newCredits, nil
}

func (s *Store) GetRooms() ([]Room, error) {
	var rooms []Room
	query := `SELECT * FROM rooms;`

	rows, err := s.db.Query(query)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var room Room
		if err := rows.Scan(&room.Id, &room.Name, &room.MaxUsers, &room.Price, &room.CreatedAt); err != nil {
			return nil, err
		}

		rooms = append(rooms, room)
	}

	return rooms, nil
}

func (s *Store) CreateRooms(rooms []models.Room) error {
	query := `INSERT INTO rooms (id, name, nb_users, price, created_at) VALUES (?, ?, ?, ?, ?)`
	for _, room := range rooms {
		_, err := s.db.Exec(query, room.Id, room.Name, room.NbUsers, room.Price, time.Now())

		if err != nil {
			return err
		}
	}

	return nil
}

func (s *Store) GetRoomByName(name string) (*models.Room, error) {
	var room models.Room

	query := `SELECT id, name, nb_users, price FROM rooms WHERE name = ? LIMIT 1;`

	err := s.db.QueryRow(query, name).Scan(
		&room.Id,
		&room.Name,
		&room.NbUsers,
		&room.Price,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("No room matching your query. \n\n Please try ./cosoft rooms to see available ones.")
		}
		return nil, err
	}

	return &room, nil

}
