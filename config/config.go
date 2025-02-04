package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type App struct {
	AppPort      string `json:"app_port"`
	AppEnv       string `json:"app_env"`
	JwtSecretKey string `json:"jwt_secret_key"`
	JwtIssuer    string `json:"jwt_issuer"`
}

type PsqlDB struct {
	Host      string `json:"host"`
	Port      string `json:"port"`
	User      string `json:"user"`
	Password  string `json:"password"`
	DBName    string `json:"dbname"`
	DBMaxOpen int    `json:"dbmaxopen"`
	DBMaxIdle int    `json:"dbmaxidle"`
}

type Config struct {
	App  App
	Psql PsqlDB
}

func parseEnvInt(key string, defaultValue int) int {
	value, err := strconv.Atoi(os.Getenv(key))
	if err != nil {
		return defaultValue
	}

	return value
}

func NewConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Printf("No .env file found. Falling back to environtment variables: %v", err)
	}

	return &Config{
		App: App{
			AppPort:      os.Getenv("APP_PORT"),
			AppEnv:       os.Getenv("APP_ENV"),
			JwtSecretKey: os.Getenv("JWT_SECRET_KEY"),
			JwtIssuer:    os.Getenv("JWT_ISSUER"),
		},
		Psql: PsqlDB{
			Host:      os.Getenv("DATABASE_HOST"),
			Port:      os.Getenv("DATABASE_PORT"),
			User:      os.Getenv("DATABASE_USER"),
			Password:  os.Getenv("DATABASE_PASSWORD"),
			DBName:    os.Getenv("DATABASE_NAME"),
			DBMaxOpen: parseEnvInt("DATABASE_MAX_OPEN", 100),
			DBMaxIdle: parseEnvInt("DATABASE_MAX_IDLE", 20),
		},
	}
}
