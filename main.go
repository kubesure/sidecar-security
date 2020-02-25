package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"

	"github.com/julienschmidt/httprouter"
	"github.com/kubesure/sidecar-security/proxy"
	"github.com/sirupsen/logrus"
)

func init() {
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logrus.InfoLevel)
}

func main() {

	router := httprouter.New()
	proxy.SetupProxy(router)

	srv := http.Server{Addr: ":8001", Handler: router}
	ctx := context.Background()
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	go func() {
		for range c {
			logrus.Info("shutting down sidecar security...")
			srv.Shutdown(ctx)
			<-ctx.Done()
		}
	}()
	logrus.Info("security sidecar started...")
	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		logrus.Fatalf("ListenAndServe(): %s", err)
	}
}
