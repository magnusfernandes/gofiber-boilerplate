package database

import (
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/magnusfernandes/gofiber-boilerplate/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var Database *gorm.DB
var err error

func InitDatabase() *gorm.DB {
	Database, err = gorm.Open(postgres.Open(getDSN()), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database")
	}

	err = Database.Exec(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp";`).Error
	if err != nil {
		log.Fatal("Error creating extensions", err)
	}
	Migrate(Database)
	return Database
}

func Migrate(db *gorm.DB) {
	err = db.AutoMigrate(
		&models.User{},
	)
	if err != nil {
		log.Fatal("Error migrating DB", err)
	}
}

func getDSN() string {
	HOST := os.Getenv("TEST_DBHOST")
	USER := os.Getenv("TEST_DBUSER")
	DBPASSWORD := os.Getenv("TEST_DBPASSWORD")
	DBNAME := os.Getenv("TEST_DBNAME")
	PORT := os.Getenv("TEST_DBPORT")
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s", HOST, USER, DBPASSWORD, DBNAME, PORT)
}
