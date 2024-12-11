package database

import (
	"fmt"
	"log"
	"sync"

	"github.com/arthurgavazza/farm-api-challenge/internal/app/infra/config"
	"github.com/arthurgavazza/farm-api-challenge/internal/app/infra/database/entities"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	db   *gorm.DB
	once sync.Once
)

func generateConnectionString(config *config.Config) string {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		config.Database.Host,
		config.Database.User,
		config.Database.Password,
		config.Database.Name,
		config.Database.Port,
	)
	return dsn

}

func NewPostgresDatabase(connectionString string) *gorm.DB {
	once.Do(func() {
		var err error
		db, err = gorm.Open(postgres.Open(connectionString), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
		})
		if err != nil {
			log.Fatalln("Failed to connect to database:", err)
		}
		db.AutoMigrate(&entities.Farm{}, &entities.CropProduction{})

	})

	return db
}
