package config

import (
	stackdriver "go.opencensus.io/exporter/stackdriver/propagation"
	"go.opencensus.io/plugin/ochttp/propagation/b3"
	oc "go.opencensus.io/trace/propagation"
	"os"
)

var (
	Port       string
	AppVersion string
	AppName    string

	Mode             string
	debug            string
	profile          string
	Base             string
	MaxOrbit         int
	tracePropagation string
	traceInternal    string
)

func Configure() {
	// General Settings
	Port = os.Getenv("PORT")
	AppName = os.Getenv("APP_NAME")
	AppVersion = os.Getenv("APP_VERSION")
	// Reactor Specific Settings
	Mode = os.Getenv("REACTOR_MODE")
	debug = os.Getenv("REACTOR_DEBUG")
	profile = os.Getenv("REACTOR_PROFILE")
	// Reactor Fixed Settings
	Base = "/rr"
	MaxOrbit = 5

	tracePropagation = os.Getenv("REACTOR_TRACE_PROPAGATION")
	traceInternal = os.Getenv("REACTOR_TRACE_INTERNAL")

}

func IsLocalMode() bool {
	return "local" == Mode
}

func IsKubernetesMode() bool {
	return "k8s" == Mode
}

func IsDebug() bool {
	return "1" == debug
}

func IsProfiling() bool {
	return "1" == profile
}

func NextOrbit() string {
	return "inf"
}

func TracePropagation() oc.HTTPFormat {
	if tracePropagation == "b3" {
		return &b3.HTTPFormat{}
	} else {
		return &stackdriver.HTTPFormat{}
	}
}
