package main

import (
	"github.com/alexvanboxel/reactor/pkg/client"
	"github.com/alexvanboxel/reactor/pkg/config"
	"github.com/alexvanboxel/reactor/pkg/reactor"
)

func main() {
	config.Configure()
	client.Init()
	reactor.Serve()
}
