// Package main is the http server of the application.
package main

import (
	"github.com/i2dou/sponge/cmd/serverNameExample_httpExample/initial"

	"github.com/i2dou/sponge/pkg/app"
)

// @title serverNameExample api docs
// @description http server api docs
// @schemes http https
// @version 2.0
// @host localhost:8080
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer your-jwt-token" to Value
func main() {
	initial.Config()
	servers := initial.RegisterServers()
	closes := initial.RegisterClose(servers)

	a := app.New(servers, closes)
	a.Run()
}
