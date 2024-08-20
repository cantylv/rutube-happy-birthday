// Copyright © ivanlobanov. All rights reserved.
package main

import (
	"github.com/cantylv/service-happy-birthday/config"
	"github.com/cantylv/service-happy-birthday/internal/app"
)

// 1day
// list of tasks
// 2) add memcache
// 8) add kafka for notification | need to add consumer and producer 
// consumer will 
// 4) write tests
func main() {
	// setup configuration
	config.Read()
	// run app
	app.Run()
}
