package resource

import (
	"go.opencensus.io/exporter/stackdriver/propagation"
	"go.opencensus.io/plugin/ochttp"
	"net/http"
)

var HttpClient *http.Client

func clientInit() {
	octr := &ochttp.Transport{
		Propagation: &propagation.HTTPFormat{},
	}
	HttpClient = &http.Client{Transport: octr}
}
