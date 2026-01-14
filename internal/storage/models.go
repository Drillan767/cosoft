package storage

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	Id          uuid.UUID `db:"id"`
	FirstName   string    `db:"first_name"`
	LastName    string    `db:"last_name"`
	Email       string    `db:"email"`
	Jwt         string    `db:"jwt"`
	Credits     int       `db:"credits"`
	ExpiresAt   time.Time `db:"expires_at"`
	SlackUserID *string   `db:"slack_user_id"`
	CreatedAt   time.Time `db:"created_at"`
}

type Room struct {
	Id        string    `db:"id"`
	Name      string    `db:"name"`
	MaxUsers  int       `db:"max_users"`
	Price     int       `db:"price"`
	CreatedAt time.Time `db:"created_at"`
}

type Reservation struct {
	Id        uint      `db:"id"`
	Date      time.Time `db:"date"`
	Room      Room      `db:"room"`
	User      User      `db:"user"`
	Duration  int       `db:"duration"`
	Cost      int       `db:"cost"`
	CreatedAt time.Time `db:"created_at"`
}
