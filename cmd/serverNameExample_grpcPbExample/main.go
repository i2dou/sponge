// Package main is the grpc server of the application.
package main

import (
	"github.com/i2dou/sponge/cmd/serverNameExample_grpcPbExample/initial"

	"github.com/i2dou/sponge/pkg/app"
)

func main() {
	initial.Config()
	servers := initial.RegisterServers()
	closes := initial.RegisterClose(servers)

	a := app.New(servers, closes)
	a.Run()
}
