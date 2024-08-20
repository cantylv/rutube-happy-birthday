// Copyright © ivanlobanov. All rights reserved.
package app

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/cantylv/service-happy-birthday/internal/crontasks"
	"github.com/cantylv/service-happy-birthday/internal/route"
	"github.com/cantylv/service-happy-birthday/internal/utils/functions"
	"github.com/cantylv/service-happy-birthday/services"
	"github.com/gorilla/mux"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// Run
// The heart of our application. Initialization of different services and server instance.
func Run() {
	logger := zap.Must(zap.NewProduction()).Sugar()
	// Tag definition for DTO.
	functions.InitValidator()
	// Initialization of DBMS, brokers, in-memory storage and etc..
	serviceCluster := services.Init()
	defer func() {
		// Shutting down all services.
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := serviceCluster.MongoClient.Disconnect(ctx); err != nil {
			logger.Errorf("error config file: %v", err)
		}
		err := serviceCluster.CacheClient.Close()
		if err != nil {
			logger.Errorf("error close memcache connections: %v", err)
		}
		err = serviceCluster.EmailBroker.Close()
		if err != nil {
			logger.Errorf("error close kafka consumer: %v", err)
		}
	}()

	// Initialization cron tasks.
	cTasks := crontasks.InitCronTasks()
	defer cTasks.Stop()

	srv := &http.Server{
		Addr:         fmt.Sprintf("%s:%d", viper.GetString("server.host"), viper.GetUint16("server.port")),
		WriteTimeout: viper.GetDuration("server.write_timeout"),
		ReadTimeout:  viper.GetDuration("server.read_timeout"),
		IdleTimeout:  viper.GetDuration("server.idle_timeout"),
		Handler: route.Initialize(
			route.RouterProps{
				R: mux.NewRouter(),
				S: serviceCluster,
			},
		),
	}

	// Run server instance.
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			logger.Fatalf("the server has terminated its work: %v", err)
		}
	}()

	// graceful shutdown
	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)

	signal.Notify(c, os.Interrupt)
	<-c

	ctx, cancel := context.WithTimeout(context.Background(), viper.GetDuration("graceful-timeout"))
	defer cancel()

	srv.Shutdown(ctx)
	logger.Info("the server has terminated its work")
	os.Exit(0)
}
