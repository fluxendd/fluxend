package models

import "time"

type Note struct {
	ID        uint      `db:"id"`
	UserId    uint      `db:"user_id"`
	Title     string    `db:"title"`
	Content   string    `db:"content"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}
