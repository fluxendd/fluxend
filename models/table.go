package models

type Table struct {
	Id     int    `db:"id"`
	Name   string `db:"name"`
	Schema string `db:"schema"`
}
