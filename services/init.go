package services

import (
	"github.com/IBM/sarama"
	"github.com/bradfitz/gomemcache/memcache"
	"github.com/cantylv/service-happy-birthday/services/kafka/consumer"
	inMemoryStorage "github.com/cantylv/service-happy-birthday/services/memcache"
	"github.com/cantylv/service-happy-birthday/services/mongodb"
	"go.mongodb.org/mongo-driver/mongo"
)

type Services struct {
	MongoClient *mongo.Client
	CacheClient *memcache.Client
	EmailBroker sarama.Consumer
}

func Init() Services {
	return Services{
		MongoClient: mongodb.Init(),
		CacheClient: inMemoryStorage.Init(),
		EmailBroker: consumer.Init(),
	}
}
