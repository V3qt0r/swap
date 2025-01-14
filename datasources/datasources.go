package datasources

import (
	"fmt"
	"log"
	"os"
	"swap/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Ds struct {
	DB *gorm.DB
}


func InitDS() (*Ds, error) {
	log.Printf("Initializing data sources\n")
	DATABASE_HOST := os.Getenv("DATABASE_HOST")
	DATABASE_USER := os.Getenv("DATABASE_USER")
	DATABASE_PASSWORD := os.Getenv("DATABASE_PASSWORD")
	DATABASE_DB := os.Getenv("DATABASE_DB")
	DATABASE_PORT := os.Getenv("DATABASE_PORT")
	DB_SSL_MODE := os.Getenv("DB_SSL_MODE")

	log.Printf("Connecting to Postgres sql\n")
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		DATABASE_HOST, DATABASE_USER, DATABASE_PASSWORD, DATABASE_DB, DATABASE_PORT, DB_SSL_MODE)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		return nil, fmt.Errorf("Error opening database: %w", err)
	}

	if err := db.AutoMigrate(
		&models.User{}, &models.Item{}, &models.SwapRequest{}, &models.Category{}, &models.Image{},
	); err != nil {
		return nil, fmt.Errorf("Error migrating models: %w", err)
	}

	return &Ds{
		DB : db,
	}, nil
}


func MigrateDB(db *gorm.DB) {
	err := db.AutoMigrate(&models.Item{}, &models.Image{})
	if err != nil {
		log.Fatalf("Migration failed: %v", err)
	}
}