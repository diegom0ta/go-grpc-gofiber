package database

import (
	"fmt"
	"log"

	"github.com/diegom0ta/go-grpc-gofiber/internal/config"
	"github.com/diegom0ta/go-grpc-gofiber/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var Db *gorm.DB

func Connect() {
	var err error

	Db, err = gorm.Open(postgres.Open(config.LoadDB()), &gorm.Config{})
	if err != nil {
		log.Printf("error connecting to database: %v", err)
	}

	log.Println("Database connected")

	err = Db.AutoMigrate(&models.User{})
	if err != nil {
		log.Printf("error running migrations: %v", err)
	}

	log.Println("Migrations done")
}

func Disconnect() error {
	instance, err := Db.DB()
	if err != nil {
		return fmt.Errorf("error while disconnecting from database: %v", err)
	}

	instance.Close()

	log.Println("Database disconnected")

	return nil
}
