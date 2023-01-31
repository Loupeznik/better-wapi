package helpers

import (
	"log"
	"os"
)

var logToFile bool

func init() {
	config := SetupIntegrationConfig()

	if config.UseLogFile {
		logToFile = true
	}
}

func Log(logLevel string, message string) {
	if logToFile {
		file := openOrCreateLogFile()
		log.SetOutput(file)
	}

	log.Printf("[%s] %s", logLevel, message)
}

func openOrCreateLogFile() *os.File {
	file, err := os.OpenFile("logs/api.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer file.Close()

	return file
}
