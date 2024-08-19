// Copyright © ivanlobanov. All rights reserved.
package main

import (
	"github.com/cantylv/service-happy-birthday/config"
	"github.com/cantylv/service-happy-birthday/internal/app"
)

// 1day
// list of tasks
// 2) add memcache
// 4) write tests
// 6) remove password from return http data
func main() {
	// setup configuration
	config.Read()
	// run app
	app.Run()
}
