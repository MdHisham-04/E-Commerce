package database

import (
	"fmt"
	"log"

	"github.com/MdHisham-04/E-Commerce/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

type Config struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}

// Connect establishes a connection to the PostgreSQL database
func Connect(config Config) error {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		config.Host, config.User, config.Password, config.DBName, config.Port,
	)

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	log.Println("Database connected successfully")
	return nil
}

// Migrate runs database migrations to create/update all tables
func Migrate() error {
	err := DB.AutoMigrate(
		&models.User{},
		&models.Product{},
		&models.CartItem{},
		&models.Order{},
		&models.OrderItem{},
	)
	if err != nil {
		return fmt.Errorf("failed to migrate database: %w", err)
	}

	log.Println("Database migrated successfully")
	return nil
}
