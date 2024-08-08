// Copyright © ivanlobanov. All rights reserved.
package app

import (
	"context"
	"time"

	"github.com/cantylv/service-happy-birthday/internal/functions"
	"github.com/cantylv/service-happy-birthday/services"
	"go.uber.org/zap"
)

func Run() {
	logger := zap.Must(zap.NewProduction()).Sugar()
	// tag definition for DTO
	functions.InitValidator()
	// initialization of DBMS, brokers, in-memory storage and etc..
	serviceCluster := services.Init()
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := serviceCluster.MongoClient.Disconnect(ctx); err != nil {
			logger.Panicf("Fatal error config file: %w.", err)
		}
	}()

	// запускаем образ сервера
	// делаем graceful shutdown
}
