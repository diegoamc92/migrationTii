package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	TiisaDBHost     string
	TiisaDBUser     string
	TiisaDBPassword string
	TiisaDBName     string
	TiisaDBPort     int

	IsbDBHost     string
	IsbDBUser     string
	IsbDBPassword string
	IsbDBName     string
	IsbDBPort     int
}

func LoadConfig() (*Config, error) {
	godotenv.Load()

	tiisaPort, _ := strconv.Atoi(getEnv("TIISA_DB_PORT", "3306"))
	isbPort, _ := strconv.Atoi(getEnv("ISB_DB_PORT", "3306"))

	return &Config{
		TiisaDBHost:     getEnv("TIISA_DB_HOST", ""),
		TiisaDBUser:     getEnv("TIISA_DB_USER", "root"),
		TiisaDBPassword: getEnv("TIISA_DB_PASSWORD", ""),
		TiisaDBName:     getEnv("TIISA_DB_NAME", ""),
		TiisaDBPort:     tiisaPort,

		IsbDBHost:     getEnv("ISB_DB_HOST", ""),
		IsbDBUser:     getEnv("ISB_DB_USER", "root"),
		IsbDBPassword: getEnv("ISB_DB_PASSWORD", ""),
		IsbDBName:     getEnv("ISB_DB_NAME", ""),
		IsbDBPort:     isbPort,
	}, nil
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
