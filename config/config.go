// Copyright © ivanlobanov. All rights reserved.
package config

import (
	"os"
	"time"

	"github.com/cantylv/service-happy-birthday/internal/utils/myconstants"
	"github.com/satori/uuid"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// setDefaultParameters
// Default configuration is develop.
func setDefaultParameters() {
	viper.AutomaticEnv()
	// common variables
	viper.SetDefault("config", "./config/dev/config.yaml")
	viper.SetDefault("graceful-timeout", 10*time.Second)

	// runtime variables
	if viper.Get(myconstants.EnvCsrfSecret) == nil {
		os.Setenv(myconstants.EnvCsrfSecret, uuid.NewV4().String())
	}
	if viper.Get(myconstants.EnvJwtSecret) == nil {
		os.Setenv(myconstants.EnvJwtSecret, uuid.NewV4().String())
	}

	// __server__ variables
	viper.SetDefault("server.host", "server")
	viper.SetDefault("server.port", 8000)
	viper.SetDefault("server.write_timeout", time.Second*5)
	viper.SetDefault("server.read_timeout", time.Second*5)
	viper.SetDefault("server.idle_timeout", time.Second*60)

	// __kafka__ variables
	viper.SetDefault("kafka.host", "kafka")
	viper.SetDefault("kafka.port", 6473)
	viper.SetDefault("kafka.topic", "emails")

	// __mongodb__ variables
	viper.SetDefault("mongo.host", "mongo")
	viper.SetDefault("mongo.port", 27017)

	// __memcache__ variables
	viper.SetDefault("memcached.host", "memcached")
	viper.SetDefault("memcached.port", 11211)
}

// getFlags
// Bind flags within current viper configuration.
func getFlags() {
	var configPath string
	var wait time.Duration

	// common flags
	pflag.StringVarP(&configPath, "config", "c", "./config/dev/config.yaml", "Defines the path to the configuration file.")
	pflag.DurationVarP(&wait, "graceful-timeout", "g", time.Second*10, "The duration for which the server gracefully wait for existing connections to finish.")
	pflag.Parse()

	// binding flags
	viper.BindPFlag("config", pflag.Lookup("config"))
	viper.BindPFlag("graceful-timeout", pflag.Lookup("graceful-timeout"))
}

// Read
// Set default configuration parameters and read config file if path is specified.
func Read() {
	logger := zap.Must(zap.NewProduction()).Sugar()
	setDefaultParameters()
	getFlags()
	viper.SetConfigFile(viper.GetString("config"))
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(*os.PathError); !ok {
			logger.Panicf("fatal error config file: %v", err)
		}
		logger.Warn("warning: configuration file is not found, programm will be executed within default configuration")
	}

	config := viper.AllSettings()
	logger.Infof("successful read of configuration, current viper configuration: %v", config)
}

// Notification
// Viper uses the following precedence order. Each item takes precedence over the item below it:
// 1) explicit call to Set
// 2) flag
// 3) env
// 4) config
// 5) key/value store
// 6) default
