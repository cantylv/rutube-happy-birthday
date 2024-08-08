package services

import (
	"github.com/cantylv/service-happy-birthday/services/mongodb"
	"go.mongodb.org/mongo-driver/mongo"
)

type Services struct {
	MongoClient *mongo.Client
}

func Init() Services {
	return Services{
		MongoClient: mongodb.Init(),
	}
}
