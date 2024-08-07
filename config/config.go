// Copyright © ivanlobanov. All rights reserved.
package config

import (
	"fmt"

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
	viper.SetDefault("nginx.port", "80")

	// __server__ variables
	viper.SetDefault("server.host", "127.0.0.1")
	viper.SetDefault("server.port", "8000")

	// __kafka__ variables
	viper.SetDefault("kafka.host", "127.0.0.1")
	viper.SetDefault("kafka.port", "6473")
}

// getFlags()
// Bind flags within current viper configuration.
func getFlags() {
	var configPath string

	pflag.StringVarP(&configPath, "config", "c", "./config/dev/config.yaml", "Defines the path to the configuration file.")
	pflag.Parse()

	// binding flags
	viper.BindPFlag("config", pflag.Lookup("config"))
}

// Read()
// Set default configuration parameters and read config file if path is specified.
func Read() {
	logger := zap.Must(zap.NewProduction()).Sugar()
	setDefaultParameters()
	getFlags()
	viper.SetConfigFile(viper.GetString("config"))
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			logger.Panicf("Fatal error config file: %w", err)
		}
		logger.Warn("Warning: configuration file is not found. Programm will be executed within default configuration.")
	}
	config := viper.AllSettings()
	fmt.Println("Current Viper Configuration:")
	for key, value := range config {
		fmt.Printf("%s: %v\n", key, value)
	}
}
