package main

import (
	"github.com/alexvanboxel/reactor/pkg/resource"
	"github.com/alexvanboxel/reactor/pkg/config"
	"github.com/alexvanboxel/reactor/pkg/reactor"
)

func main() {
	config.Configure()
	resource.Init()
	reactor.Serve()
}
