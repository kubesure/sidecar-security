package middleware

import (
	"net/http"
	"net/http/httputil"
	"os"

	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)
}

//Logger middleware function logs request params
func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Infof("logging: request middleware")
		next.ServeHTTP(w, r)
		log.Infof("logging: response middleware")
	})
}

//Auth middleware function authenticates request
func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Infof("Auth: Request before proxy call")
		if r.Header.Get("user") != "foo" {
			w.WriteHeader(http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
		log.Infof("Auth: Pass")
	})
}

//Final middleware
func Final(proxy *httputil.ReverseProxy) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("Passing call to origin")
		proxy.ServeHTTP(w, r)
	})
}
