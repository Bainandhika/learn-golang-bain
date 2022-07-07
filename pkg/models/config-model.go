package models

type Config struct {
	DB DBConfig
	Logger LoggerConfig
}

type DBConfig struct {
	DBDriver	string
	Host		string
	Port		string
	UserName	string
	Password	string
	DBName		string
}

type LoggerConfig struct {
	LogPath string
}