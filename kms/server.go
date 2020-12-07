package main

import (
	"net/http"
	"os"

	"go-samples/kms/kmsapi"
	"go-samples/kms/routes"

	log "github.com/sirupsen/logrus"
)

var (
	fileHandle  *os.File
	port        string
	logFilePath string
)

func main() {

	// KMS log file path
	logFilePath := os.Getenv("KMS_LOG_PATH_FILENAME")
	if logFilePath == "" {
		log.Fatalln("Error: missing KMS service Log path and filename")
	}

	if logFilePath != "" {
		var err error
		if fileHandle, err = os.OpenFile(logFilePath, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644); err != nil {
			log.SetOutput(os.Stderr)
			log.SetLevel(log.FatalLevel)
			log.WithError(err).Fatalf("Error, could not open log file to log: %q", logFilePath)
		}
	}
	log.SetOutput(fileHandle)
	defer fileHandle.Close()

	// Set the debug log to Info level
	//log.SetLevel(log.InfoLevel)
	log.SetLevel(log.DebugLevel)

	// kms initialise
	if err := kmsapi.Init(); err != nil {
		log.Fatalf("Error: KMS Init: %+v\n", err)
	}

	// service port init
	port := os.Getenv("KMS_PORT")
	if port == "" {
		log.Fatalln("Error: missing KMS service rest port")
	}

	// Handle routes
	http.Handle("/", routes.RHandlers())

	// KMS Service
	log.Infoln("KMS service up on port :", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
