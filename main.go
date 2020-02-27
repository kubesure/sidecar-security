package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"

	"github.com/julienschmidt/httprouter"
	"github.com/kubesure/sidecar-security/proxy"
	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)
}

func main() {

	router := httprouter.New()
	proxy.SetupProxy(router)

	srv := http.Server{Addr: ":8000", Handler: router}
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
