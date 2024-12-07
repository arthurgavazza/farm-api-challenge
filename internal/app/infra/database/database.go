package database

import (
	"fmt"
	"log"
	"sync"

	"github.com/arthurgavazza/farm-api-challenge/internal/app/infra/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	db   *gorm.DB
	once sync.Once
)

// NewPostgresDatabase initializes a new GORM database connection.
func NewPostgresDatabase(config *config.Config) *gorm.DB {
	once.Do(func() {
		dsn := fmt.Sprintf(
			"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
			config.Database.Host,
			config.Database.User,
			config.Database.Password,
			config.Database.Name,
			config.Database.Port,
		)

		var err error
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Info), 
		})
		if err != nil {
			log.Fatalln("Failed to connect to database:", err)
		}
		db.AutoMigrate()
	   
	})

	return db
}
