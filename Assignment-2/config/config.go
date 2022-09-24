package config

import (
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/joho/godotenv"
)

type Config struct {
	ServerAddress string
	Database      Database
	BaseUrl       string

	JWT_EXP            time.Duration
	JWT_SIGNING_METHOD jwt.SigningMethod
	JWT_SECRET_KEY     string
}

type Database struct {
	Username string
	Password string
	Address  string
	Port     string
	Name     string
}

var config Config

func Init() {
	err := godotenv.Load()
	if err != nil {
		log.Printf("ERROR .env Not found")
	}

	config.ServerAddress = os.Getenv("SERVER_ADDRESS")
	config.Database.Username = os.Getenv("DB_USERNAME")
	config.Database.Password = os.Getenv("DB_PASSWORD")
	config.Database.Address = os.Getenv("DB_ADDRESS")
	config.Database.Port = os.Getenv("DB_PORT")
	config.Database.Name = os.Getenv("DB_NAME")

	config.JWT_EXP = time.Duration(120) * time.Hour
	config.JWT_SIGNING_METHOD = jwt.SigningMethodHS256
	config.JWT_SECRET_KEY = os.Getenv("JWT_SCRET_KEY")
}

func GetConfig() *Config {
	return &config
}
