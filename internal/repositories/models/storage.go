package models

import "time"

type Storage struct {
	ID          int        `db:"id"`
	Name        string     `db:"name"`
	Temperature float32    `db:"temperature"`
	Humidity    float32    `db:"humidity"`
	CreatedAt   *time.Time `db:"created_at"`
	Type        Type
}

type Type struct {
	ID   int    `db:"id,id_type"`
	Name string `db:"name"`
}
