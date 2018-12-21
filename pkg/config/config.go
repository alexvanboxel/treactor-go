package config

import "os"

var (
	Port     string
	Version  string
	Name     string
	Mode     string
	Base     string
	MaxOrbit int
)

func Configure() {
	Port = os.Getenv("PORT")
	Version = os.Getenv("VERSION")
	Name = os.Getenv("NAME")
	Mode = os.Getenv("MODE")
	Base = "/rr"
	MaxOrbit = 5
}

func IsLocal() bool {
	return "local" == Mode
}

func NextOrbit() string {
	return "inf"
}
