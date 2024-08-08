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
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	connLine := fmt.Sprintf("mongodb://%s:%d", viper.GetString("mongodb.host"), viper.GetUint16("mongodb.port"))
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(connLine))
	if err != nil {
		logger.Panicf("Fatal error MongoDB connect: %w.", err)
	}
	return client
}
