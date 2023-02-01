package helpers

import (
	"errors"
	"fmt"
	"io"
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
		file := *openOrCreateLogFile()
		defer file.Close()

		wrt := io.MultiWriter(os.Stdout, &file)
		log.SetOutput(wrt)
	}

	log.Printf("[%s] %s", logLevel, message)
}

func openOrCreateLogFile() *os.File {
	dir := "logs"

	if _, err := os.Stat(dir); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(dir, os.ModePerm)
		if err != nil {
			log.Fatalf("error creating log directory: %v", err)
		}
	}

	file, err := os.OpenFile(fmt.Sprintf("%s/api.log", dir), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}

	return file
}
