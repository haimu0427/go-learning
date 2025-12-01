package contact

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func loadEnv() error {
	err := godotenv.Load()
	if err != nil {
		return fmt.Errorf("error loading .env file: %v", err)
	}
	return nil
}
func LoadEnvlog() error {
	// Load environment variables
	if err := loadEnv(); err != nil {
		log.Fatal(err)
	}
	return nil
}

func BuildClientURL() string {
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	return fmt.Sprintf("redis://%s:%s@%s:%s", dbUser, dbPassword, dbHost, dbPort)
}
