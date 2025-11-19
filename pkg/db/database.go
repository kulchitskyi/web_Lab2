package db

import (
	"log"

	//"deu/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDB(dsn string) *gorm.DB {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	//err = db.AutoMigrate(
	//	&models.User{},
	//	&models.Place{},
	//	&models.UserPlace{},
	//)
	//if err != nil {
	//	log.Fatalf("Failed to auto-migrate database: %v", err)
	//}

	log.Println("Database connection established and migrations complete.")
	return db
}