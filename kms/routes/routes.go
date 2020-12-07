package routes

import (
	"net/http"
	"os"
	"strings"

	"go-samples/kms/handlers"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

// RHandlers of routers.
func RHandlers() *mux.Router {
	r := mux.NewRouter().StrictSlash(true)
	r.Use(CommonMiddleware)
	r.HandleFunc("/", handlers.Health).Methods("GET")
	r.HandleFunc("/kms", handlers.RequestKMS).Methods("POST")
	return r
}

// CommonMiddleware --Set content-type
func CommonMiddleware(next http.Handler) http.Handler {
	// set logging level based on environment var
	logLevel := os.Getenv("KMS_LOGGING_LEVEL")
	switch strings.ToLower(strings.Trim(logLevel, "")) {
	case "info":
		log.SetLevel(log.InfoLevel)
	case "debug":
		log.SetLevel(log.DebugLevel)
	default:
		log.Warnln(" invalid log level, except only info/debug")
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, Access-Control-Request-Headers, Access-Control-Request-Method, Connection, Host, Origin, User-Agent, Referer, Cache-Control, X-header")
		next.ServeHTTP(w, r)
	})
}
