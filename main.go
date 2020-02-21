package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetLevel(log.DebugLevel)
	log.SetOutput(os.Stdout)
}

func main() {
	log.Info("sidecar security starting...")
	mux := http.NewServeMux()
	mux.HandleFunc("/", invoke)
	srv := http.Server{Addr: ":8000", Handler: mux}
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

func invoke(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	data := (time.Now()).String()
	log.Debug("Invoked....")
	w.Write([]byte(data))
}
