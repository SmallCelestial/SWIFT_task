package db

import (
	"SWIFT_task/internal/model"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

func InitDB() *gorm.DB {
	//host := os.Getenv("DB_HOST")
	//user := os.Getenv("DB_USER")
	//password := os.Getenv("DB_PASSWORD")
	//dbname := os.Getenv("DB_NAME")
	//port := "5432"
	host := "localhost"
	user := "postgres"
	password := "admin"
	dbname := "swift_task"
	port := "5432"
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", host, user, password, dbname, port)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	err = db.AutoMigrate(&model.Bank{}, &model.BankRelationship{})
	if err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	return db
}
