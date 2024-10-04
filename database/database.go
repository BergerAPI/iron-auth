package database

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
)

var (
	Instance *gorm.DB
)

func Init(url string) {
	db, err := gorm.Open(sqlite.Open(url), &gorm.Config{})

	if err != nil {
		log.Fatalln(err)
	}

	_ = db.AutoMigrate(&User{})
	_ = db.AutoMigrate(&Client{})

	Instance = db
}
