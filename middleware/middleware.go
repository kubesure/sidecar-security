package middleware

import (
	"net/http"
	"net/http/httputil"
	"os"

	"github.com/sirupsen/logrus"
)

func init() {
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logrus.InfoLevel)
}

//Log middleware type logs the request
type Log struct {
	handler http.Handler
	log     string
}

//Authenticate middleware authenticates the user
type Authenticate struct {
	handler http.Handler
}

//NewLogger middleware type logs the request
func NewLogger(handler http.Handler, log string) *Log {
	return &Log{handler: handler, log: log}
}

//Handler for middleware
type Handler interface {
	ServeHTTP(http.ResponseWriter, *http.Request)
}

//Logger Middleware
func (logger *Log) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	logrus.Infof("request -%s", "request middleware")
	logger.ServeHTTP(w, r)
}

func (auth *Authenticate) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	logrus.Infof("Auth -%s", "ok...")
	auth.ServeHTTP(w, r)
}

//Logger middleware function logs request params
func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logrus.Infof("logging: request middleware")
		next.ServeHTTP(w, r)
		logrus.Infof("logging: response middleware")
	})
}

//Auth middleware function authenticates request
func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logrus.Infof("Auth: Request before proxy call")
		if r.Header.Get("user") != "foo" {
			w.WriteHeader(http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
		logrus.Infof("Auth: Pass")
	})
}

//Final middleware
func Final(proxy *httputil.ReverseProxy) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logrus.Println("Passing call to origin")
		proxy.ServeHTTP(w, r)
	})
}
