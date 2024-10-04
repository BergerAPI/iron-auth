package database

import "time"

type User struct {
	Id        string `gorm:"primaryKey"`
	Email     string
	Password  string
	CreatedAt time.Time
}

type Client struct {
	Id          string `gorm:"primaryKey"`
	Name        string
	RedirectUri string
	Secret      string
	CreatedAt   time.Time
}
