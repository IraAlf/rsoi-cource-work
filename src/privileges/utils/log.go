package utils

import (
	"log"
	"net/http"
	"os"

	"github.com/urfave/negroni"
)

var (
	Logger  *log.Logger
	logFile *os.File
)

func InitLogger(filename ...string) {
	path := Config.LogFile

	if len(filename) != 0 {
		path = filename[0]
	}

	logFile, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		log.Print(err)
	}

	Logger = log.New(logFile, "", log.Ldate|log.Ltime)
}

func CloseLogger() {
	logFile.Close()
}

var LogHandler = func(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		Logger.Printf("%s REQUEST\t URL:%s \tAddress: %s", r.Method, r.URL, r.RemoteAddr)
		log.Printf("%s REQUEST\t URL:%s \tAddress: %s\n", r.Method, r.URL, r.RemoteAddr)

		lrw := negroni.NewResponseWriter(w)
		next.ServeHTTP(lrw, r)

		statusCode := lrw.Status()
		log.Printf("<-- %d %s\n", statusCode, http.StatusText(statusCode))
	})
}
