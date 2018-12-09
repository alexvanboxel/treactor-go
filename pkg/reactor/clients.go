package reactor

import (
	"go.opencensus.io/exporter/stackdriver/propagation"
	"go.opencensus.io/plugin/ochttp"
	"net/http"
)

var HttpClient *http.Client

func Init() {

	octr := &ochttp.Transport{
		Propagation:&propagation.HTTPFormat{},
	}

	//tr := &http.Transport{
	//	MaxIdleConns:       10,
	//	IdleConnTimeout:    30 * time.Second,
	//	DisableCompression: true,
	//}
	HttpClient = &http.Client{Transport: octr}
}
