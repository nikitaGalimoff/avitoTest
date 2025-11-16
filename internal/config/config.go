package config

import (
	"fmt"
	"os"
)

type Config struct {
	ServerPort   string
	DBHost       string
	DBPort       string
	DBUser       string
	DBPassword   string
	DBName       string
	AdminToken   string
	UserToken    string
	MigratorPath string
}

func Load() *Config {
	return &Config{
		ServerPort:   getEnv("SERVER_PORT", "8080"),
		DBHost:       getEnv("DB_HOST", "postgres"),
		DBPort:       getEnv("DB_PORT", "5432"),
		DBUser:       getEnv("DB_USER", "postgres"),
		DBPassword:   getEnv("DB_PASSWORD", "postgres"),
		DBName:       getEnv("DB_NAME", "avitotest"),
		AdminToken:   getEnv("ADMIN_TOKEN", "admin-token"),
		UserToken:    getEnv("USER_TOKEN", "user-token"),
		MigratorPath: getEnv("/migrations", "migrations-path"),
	}
}

func (c *Config) GetDBConnectionString() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		c.DBHost, c.DBPort, c.DBUser, c.DBPassword, c.DBName)
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
