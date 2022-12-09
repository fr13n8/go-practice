package database

import (
	"github.com/fr13n8/go-practice/internal/config"
	"github.com/fr13n8/go-practice/internal/domain"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewDb(cfg *config.DatabaseConfig) *gorm.DB {
	dsn := "host=" + cfg.Host + " user=" + cfg.User + " password=" + cfg.Password + " port=" + cfg.Port + " sslmode=disable TimeZone=Europe/Moscow"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	err = db.Table(cfg.Name).AutoMigrate(&domain.Task{})
	if err != nil {
		return nil
	}

	db = db.Table(cfg.Name)
	return db
}
