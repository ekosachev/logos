package database

import (
	"fmt"

	"github.com/ekosachev/logos/internal/config"
	"github.com/ekosachev/logos/internal/repositories"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectToDb() (*gorm.DB, error) {
	cfg := config.GetConfig()
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC",
		cfg.DBHost,
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBName,
		cfg.DBPort,
	)

	db, err := gorm.Open(postgres.Open(dsn))
	if err != nil {
		return nil, err
	}

	if err = db.AutoMigrate(&repositories.User{}); err != nil {
		return nil, err
	}

	return db, nil
}
