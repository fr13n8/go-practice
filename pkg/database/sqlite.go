package database

import (
	"github.com/fr13n8/go-practice/internal/domain"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func NewDb(database string) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(database), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	err = db.Table("tasks").AutoMigrate(&domain.Task{})
	if err != nil {
		return nil
	}
	return db
}
