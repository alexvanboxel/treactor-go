package config

import "os"

var (
	Port    string
	Version string
	Name    string
	Mode    string
)

func Configure() {
	Port = os.Getenv("PORT")
	Version = os.Getenv("VERSION")
	Name = os.Getenv("NAME")
	Mode = os.Getenv("MODE")
}
