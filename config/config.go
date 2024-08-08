// Copyright © ivanlobanov. All rights reserved.
package config

import (
	"os"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// setDefaultParameters()
// Default configuration is develop.
func setDefaultParameters() {
	// common variables
	viper.SetDefault("config", "./config/dev/config.yaml")

	// __nginx__ variables
	viper.SetDefault("nginx.host", "127.0.0.1")
	viper.SetDefault("nginx.port", 80)

	// __server__ variables
	viper.SetDefault("server.host", "127.0.0.1")
	viper.SetDefault("server.port", 8000)

	// __kafka__ variables
	viper.SetDefault("kafka.host", "127.0.0.1")
	viper.SetDefault("kafka.port", 6473)

	// __mongodb__ variables
	viper.SetDefault("mongodb.host", "127.0.0.1")
	viper.SetDefault("mongodb.port", 27017)
}

// getFlags()
// Bind flags within current viper configuration.
func getFlags() {
	var configPath string
	var host string

	// common flags
	pflag.StringVarP(&configPath, "config", "c", "./config/dev/config.yaml", "Defines the path to the configuration file.")
	pflag.StringVarP(&host, "host", "h", "127.0.0.1", "Defines the ip-address of the host.")
	pflag.Parse()

	// binding flags
	viper.BindPFlag("config", pflag.Lookup("config"))

	viper.BindPFlag("nginx.host", pflag.Lookup("host"))
	viper.BindPFlag("server.host", pflag.Lookup("host"))
	viper.BindPFlag("kafka.host", pflag.Lookup("host"))
	viper.BindPFlag("mongodb.host", pflag.Lookup("host"))
}

// Read()
// Set default configuration parameters and read config file if path is specified.
func Read() {
	logger := zap.Must(zap.NewProduction()).Sugar()
	setDefaultParameters()
	getFlags()
	viper.SetConfigFile(viper.GetString("config"))
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(*os.PathError); !ok {
			logger.Panicf("Fatal error config file: %v.", err)
		}
		logger.Warn("Warning: configuration file is not found. Programm will be executed within default configuration.")
	}

	config := viper.AllSettings()
	logger.Infof("Current Viper Configuration: %v", config)
}

// Notification
// Viper uses the following precedence order. Each item takes precedence over the item below it:
// 1) explicit call to Set
// 2) flag
// 3) env
// 4) config
// 5) key/value store
// 6) default
