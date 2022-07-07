package configs

import (
	"learn-golang-bain/pkg/models"
	"log"
	"os"
	"regexp"

	"github.com/joho/godotenv"
)

const projectDirName = "learn-golang-bain" // change to relevant project name

func loadEnv() {
	projectName := regexp.MustCompile(`^(.*` + projectDirName + `)`)
	currentWorkDirectory, _ := os.Getwd()
	rootPath := projectName.Find([]byte(currentWorkDirectory))

	err := godotenv.Load(string(rootPath) + `\configs\config.env`)

	if err != nil {
		log.Fatalf("Error loading .env file. Detail : %s", err.Error())
	}
}

func GetConfig() models.Config {
	loadEnv()

	dbDriver := os.Getenv("DB_DRIVER")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_DATABASE")
	dbUsername := os.Getenv("DB_USERNAME")
	dbPassword := os.Getenv("DB_PASSWORD")

	logPath := os.Getenv("LOG_PATH")

	return models.Config{
		DB: models.DBConfig{
			DBDriver:	dbDriver,
			Host:       dbHost,
			Port:       dbPort,
			UserName:   dbUsername,
			Password:   dbPassword,
			DBName:		dbName,
		},
		Logger: models.LoggerConfig{
			LogPath: logPath,
		},
	}
}