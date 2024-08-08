// Copyright © ivanlobanov. All rights reserved.
package main

import (
	"github.com/cantylv/service-happy-birthday/config"
	"github.com/cantylv/service-happy-birthday/internal/app"
)

// i want to use
// viper, mongoDB

func main() {
	// чтение конфигурации
	config.Read()
	// запуск приложения
	app.Run()
}
