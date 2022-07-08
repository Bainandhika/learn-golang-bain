package models

type Config struct {
	DB		DBConfig
	Logger	LoggerConfig
	Redis	RedisConfig
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

type RedisConfig struct {
	Host			string
	Port			string
	User			string
	Password		string
	SentinelIP1		string
	SentinelIP2		string
	SentinelIP3		string
	SentinelPort	string
	MasterName		string
	AuthHA			string
	EnableHA		string
}