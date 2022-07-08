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

	logPath := os.Getenv("LOG_PATH")

	dbDriver := os.Getenv("DB_DRIVER")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_DATABASE")
	dbUsername := os.Getenv("DB_USERNAME")
	dbPassword := os.Getenv("DB_PASSWORD")

	redisHost := os.Getenv("REDIS_HOST")
	redisPort := os.Getenv("REDIS_PORT")
	redisUser := os.Getenv("REDIS_USER")
	redisPassword := os.Getenv("REDIS_PASSWORD")
	redisSentinelIP1 := os.Getenv("REDIS_SENTINEL_IP1")
	redisSentinelIP2 := os.Getenv("REDIS_SENTINEL_IP2")
	redisSentinelIP3 := os.Getenv("REDIS_SENTINEL_IP3")
	redisSentinelPort := os.Getenv("REDIS_SENTINEL_PORT")
	redisMasterName := os.Getenv("REDIS_MASTERNAME")
	enableHA := os.Getenv("REDIS_HA")

	return models.Config{
		Logger: models.LoggerConfig{
			LogPath: logPath,
		},
		DB: models.DBConfig{
			DBDriver:	dbDriver,
			Host:       dbHost,
			Port:       dbPort,
			UserName:   dbUsername,
			Password:   dbPassword,
			DBName:		dbName,
		},
		Redis: models.RedisConfig{
			Host:			redisHost,
			Port:			redisPort,
			User:			redisUser,
			Password:		redisPassword,
			SentinelIP1: 	redisSentinelIP1,
			SentinelIP2:	redisSentinelIP2,
			SentinelIP3:	redisSentinelIP3,
			SentinelPort:	redisSentinelPort,
			MasterName:		redisMasterName,
			EnableHA:		enableHA,
		},
	}
}