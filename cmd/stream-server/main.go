package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	runtime "github.com/ahab94/streaming-service"
	"github.com/ahab94/streaming-service/api"
	"github.com/ahab94/streaming-service/config"
	"github.com/ahab94/streaming-service/streaming"
)

var logger *logrus.Logger

func main() {
	ctx := context.TODO()

	logger = logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})

	// instantiate runtime
	rt, err := runtime.NewRuntime()
	if err != nil {
		panic(err)
	}

	// configure stream server
	streamSrv, err := streaming.NewStreamServer(ctx, rt.Service(),
		rt.Publisher().AddUserInput(), rt.Publisher().RemoveUserInput())
	if err != nil {
		panic(err)
	}

	// configure http server
	httpSrv := api.NewStreamService(ctx, rt.Service())

	httpSrv.ReadTimeout = 5 * time.Second
	httpSrv.WriteTimeout = 10 * time.Second
	httpSrv.Addr = fmt.Sprintf("%s:%s", viper.GetString(config.ServerHost),
		viper.GetString(config.ServerPort))

	// graceful shutdown setup
	done := make(chan bool, 1)
	quit := make(chan os.Signal, 1)

	signal.Notify(quit, os.Interrupt)

	go graceful(ctx, httpSrv, streamSrv, rt.Publisher(), quit, done)

	// start reader, publisher and http and stream server
	go rt.Publisher().Start()

	logger.Infof("Server is ready to handle tcp requests at %s", streamSrv.Addr())
	go streamSrv.ListenAndServe()

	go rt.Reader().Start()

	logger.Infof("Server is ready to handle requests at %+v", httpSrv.Addr)
	if err := httpSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Fatalf("Could not listen on %s: %v", httpSrv.Addr, err)
	}

	<-done
	logger.Info("Server stopped gracefully")
}

func graceful(ctx context.Context, server *http.Server, streamSvr *streaming.StreamServer, publisher *streaming.Publisher, quit <-chan os.Signal, done chan<- bool) {
	<-quit

	logger.Info("server is shutting down...")

	server.SetKeepAlivesEnabled(false)

	if err := server.Shutdown(ctx); err != nil {
		logger.Fatalf("could not gracefully graceful the server: %+v", err)
	}

	streamSvr.Stop()
	publisher.Stop()
	close(done)
}
