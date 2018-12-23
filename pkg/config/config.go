package config

import "os"

var (
	Port       string
	AppVersion string
	AppName    string

	Mode     string
	Base     string
	MaxOrbit int
)

func Configure() {
	// General Settings
	Port = os.Getenv("PORT")
	AppName = os.Getenv("APP_NAME")
	AppVersion = os.Getenv("APP_VERSION")
	// Reactor Specific Settings
	Mode = os.Getenv("REACTOR_MODE")
	Mode = os.Getenv("REACTOR_DEBUG")
	// Reactor Fixed Settings
	Base = "/rr"
	MaxOrbit = 5
}

func IsLocal() bool {
	return "local" == Mode
}

func NextOrbit() string {
	return "inf"
}
