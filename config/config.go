package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	DBHost     string
	DBUser     string
	DBPassword string
	DBName     string
	DBPort     int
}

func LoadConfig() (*Config, error) {
	godotenv.Load()
	port, _ := strconv.Atoi(getEnv("DB_PORT", "3306"))

	return &Config{
		DBHost:     getEnv("DB_HOST", "dev-vulcanomigracion.cy5rvtfr1tpj.us-east-1.rds.amazonaws.com"),
		DBUser:     getEnv("DB_USER", "root"),
		DBPassword: getEnv("DB_PASSWORD", "LLPmRkMxAVTCaPbauPQJ"),
		DBName:     getEnv("DB_NAME", "isb"),
		DBPort:     port,
	}, nil
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
