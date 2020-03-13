package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/kubesure/sidecar-security/proxy"
	log "github.com/sirupsen/logrus"
)

//Initializes logurs with info level
func init() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)
}

//Initializes a http request router and starts the sidecar and listens port 8000
//Shuts down the sidecar gracefully on a os interrupt.
func main() {

	router := httprouter.New()
	proxy.SetupProxy(router)

	srv := http.Server{
		Addr:         ":8000",
		Handler:      http.TimeoutHandler(router, 1*time.Second, "Timeout!!"),
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 2 * time.Second,
	}
	ctx := context.Background()
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	go func() {
		for range c {
			log.Info("shutting down sidecar security...")
			srv.Shutdown(ctx)
			<-ctx.Done()
		}
	}()
	log.Info("security sidecar started...")
	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("ListenAndServe(): %s", err)
	}
}
