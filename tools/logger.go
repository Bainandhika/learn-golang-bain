package tools

import (
	"fmt"
	"learn-golang-bain/configs"
	"log"
	"os"
	"time"
)

var (
	LogWarning *log.Logger
	LogPanic   *log.Logger
	LogInfo    *log.Logger
	LogDebug   *log.Logger
	LogError   *log.Logger
	LogFatal   *log.Logger
)

func Logger() {
	logPath := configs.GetConfig().Logger.LogPath
	timeNow := time.Now()
	logFileName := fmt.Sprintf("applog-%d%02d%02d.log", timeNow.Year(), timeNow.Month(), timeNow.Day())

	var logDirectory string
	if logPath == "" {
		logDirectory = logFileName
	} else {
		logDirectory = logPath+"/"+logFileName
	}

	LogFile, err := os.OpenFile(logDirectory, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}

	LogInfo = log.New(LogFile, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	LogWarning = log.New(LogFile, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
	LogFatal = log.New(LogFile, "FATAL: ", log.Ldate|log.Ltime|log.Lshortfile)
	LogPanic = log.New(LogFile, "PANIC: ", log.Ldate|log.Ltime|log.Lshortfile)
	LogError = log.New(LogFile, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
	LogDebug = log.New(LogFile, "DEBUG: ", log.Ldate|log.Ltime|log.Lshortfile)
}
