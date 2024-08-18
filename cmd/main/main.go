// Copyright © ivanlobanov. All rights reserved.
package main

import (
	"github.com/cantylv/service-happy-birthday/config"
	"github.com/cantylv/service-happy-birthday/internal/app"
)

// list of tasks
// 1) add hash and salt for password
// 2) add memcache
// 3) add docker-compose.yaml
// 4) write tests
// 5) add check that user can't change interval if he is not subscribed on employee
func main() {
	// setup configuration
	config.Read()
	// run app
	app.Run()
}
