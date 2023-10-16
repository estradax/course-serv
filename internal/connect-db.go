package internal

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func ConnectDB() (*gorm.DB, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, err
	}

	PG_USER := os.Getenv("POSTGRES_USER")
	PG_PASSWORD := os.Getenv("POSTGRES_PASSWORD")
	PG_DB := os.Getenv("POSTGRES_DB")

	dsn := fmt.Sprintf("host=localhost user=%v password=%v dbname=%v port=5432 sslmode=disable", PG_USER, PG_PASSWORD, PG_DB)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, err
	}

	return db, nil
}
