package resource

import (
	"github.com/alexvanboxel/reactor/pkg/config"
	"go.opencensus.io/plugin/ochttp"
	"net/http"
)

var HttpClient *http.Client

func clientInit() {
	octr := &ochttp.Transport{
		Propagation: config.TracePropagation(),
	}
	HttpClient = &http.Client{Transport: octr}
}
