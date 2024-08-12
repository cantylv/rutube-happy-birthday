package services

import (
	"github.com/bradfitz/gomemcache/memcache"
	inMemoryStorage "github.com/cantylv/service-happy-birthday/services/memcache"
	"github.com/cantylv/service-happy-birthday/services/mongodb"
	"go.mongodb.org/mongo-driver/mongo"
)

type Services struct {
	MongoClient *mongo.Client
	CacheClient *memcache.Client
}

func Init() Services {
	return Services{
		MongoClient: mongodb.Init(),
		CacheClient: inMemoryStorage.Init(),
	}
}
