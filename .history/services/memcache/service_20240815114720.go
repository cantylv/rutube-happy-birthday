package memcache

import (
	"fmt"
	"time"

	"github.com/bradfitz/gomemcache/memcache"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func Init() *memcache.Client {
	logger := zap.Must(zap.NewProduction()).Sugar()

	connLine := fmt.Sprintf("%s:%d", viper.GetString("memcache.host"), viper.GetUint16("memcache.port"))
	client := memcache.New(connLine)
	client.MaxIdleConns = 20
	client.Timeout = 5 * time.Second

	logger.Info("Succesful connection to Memcache.")
	return client
}
