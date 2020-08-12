package main

import (
	"github.com/alexvanboxel/treactor-go/pkg/resource"
	"github.com/alexvanboxel/treactor-go/pkg/config"
	"github.com/alexvanboxel/treactor-go/pkg/reactor"
)

func main() {
	config.Configure()
	resource.Init()
	reactor.Serve()
}
