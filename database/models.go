package database

import "time"

type User struct {
	Id        string `gorm:"primaryKey"`
	Email     string
	Password  string
	CreatedAt time.Time
}
