// Copyright © ivanlobanov. All rights reserved.
package memcache

import (
	"fmt"
	"time"

	"github.com/bradfitz/gomemcache/memcache"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// I want to use memcache, but it doesn't fit out data model.
func Init() *memcache.Client {
	logger := zap.Must(zap.NewProduction()).Sugar()

	connLine := fmt.Sprintf("%s:%d", viper.GetString("memcached.host"), viper.GetUint16("memcached.port"))
	client := memcache.New(connLine)
	client.MaxIdleConns = 20
	client.Timeout = 5 * time.Second

	logger.Info("succesful connection to Memcached")
	return client
}
