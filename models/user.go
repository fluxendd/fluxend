package models

import (
	"fmt"
	"time"
)

const UserStatusActive = "active"
const UserStatusInactive = "inactive"

type User struct {
	ID        uint      `db:"id"`
	Username  string    `db:"username"`
	Email     string    `db:"email"`
	Status    string    `db:"status"`
	Bio       string    `db:"bio"`
	Password  string    `db:"password"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

func (u User) GetFields() string {
	return fmt.Sprintf("id, username, email, status, bio, password, created_at, updated_at")
}

func (u User) IsActive() bool {
	return u.Status == UserStatusActive
}
