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

func createMongoCollection(name string) {
	options := options.CreateCollection().
}

func Init() *mongo.Client {
	logger := zap.Must(zap.NewProduction()).Sugar()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	connLine := fmt.Sprintf(`mongodb://%s:%s@%s:%d/main?
		directConneciton=true&tls=false&
		connectTimeoutMS=5000&socketTimeoutMS=10000&
		maxPoolSize=30&minPoolSize=0&maxConnecting=3`,
		viper.GetString("mongodb.root"), viper.GetString("mongodb.root_password"),
		viper.GetString("mongodb.host"), viper.GetUint16("mongodb.port"))
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(connLine))
	if err != nil {
		logger.Panicf("Fatal error MongoDB connect: %w.", err)
	}
	i := 0
	for client.Ping(ctx, nil) != nil {
		logger.Errorf("Error MongoDB ping: %w.", err)
		time.Sleep(2 * time.Second)
		i++
		if i == 3 {
			logger.Panicf("Fatal error MongoDB ping: %w.", err)
		}
	}

	// creating collections
	createMongoCollection("user")
	createMongoCollection("subs")
	logger.Info("Succesful connection to MongoDB.")
	return client
}
