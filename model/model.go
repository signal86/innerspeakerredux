package model

import (
	"encoding/hex"
	"fmt"
	"os"

	"crypto/sha256"

	"github.com/joho/godotenv"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
	if err := godotenv.Load(); err != nil {
		panic("no env file :brokenheart:")
	}

	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/innerspeaker_storage?parseTime=True",
		user, password, host, port)

	database, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	database.AutoMigrate(&User{})

	DB = database
}

// Tuff sha256 hashing
func Hash(password string) string {
	hasher := sha256.New()
	hasher.Write([]byte(password))
	hashed := hex.EncodeToString(hasher.Sum(nil))
	return string(hashed)
}
