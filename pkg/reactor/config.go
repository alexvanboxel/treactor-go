package reactor

type Configuration struct {
	Version string
}

var Config *Configuration

func Configure() {
	Config = &Configuration{}
}
