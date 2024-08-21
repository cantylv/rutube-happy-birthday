// Copyright © ivanlobanov. All rights reserved.
package mongodb

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

func Init() *mongo.Client {
	logger := zap.Must(zap.NewProduction()).Sugar()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	connLine := fmt.Sprintf(`mongodb://%s:%d/main?connectTimeoutMS=5000&socketTimeoutMS=10000&maxPoolSize=30&minPoolSize=0&maxConnecting=3`,
		viper.GetString("mongo.host"), viper.GetUint16("mongo.port"))
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(connLine))
	if err != nil {
		logger.Panicf("fatal error MongoDB connect: %w", err)
	}
	for i := 0; i < 3; i++ {
		err = client.Ping(ctx, nil)
		if err == nil {
			break
		}
		logger.Infof("mongoDB ping n%d: %v", i, err)
		time.Sleep(3 * time.Second)
		if i == 2 {
			logger.Panicf("fatal error MongoDB ping: %v", err)
		}
	}

	logger.Info("succesful connection to MongoDB")
	return client
}
