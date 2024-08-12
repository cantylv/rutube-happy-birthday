// Copyright © ivanlobanov. All rights reserved.
package main

import (
	"github.com/cantylv/service-happy-birthday/config"
	"github.com/cantylv/service-happy-birthday/internal/app"
)

func main() {
	// setup configuration
	config.Read()
	// run app
	app.Run()
}
