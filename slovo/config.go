package slovo

import "net/http"

type Config struct {
	Debug      bool
	ConfigFile string
	Serve      ServeConfig
	ServeCGI   ServeCGIConfig
}

type ServeConfig struct {
	Location string
}

// Minimum ENV values for emulating a CGI request on the command line
type ServeCGIConfig struct {
	REQUEST_METHOD  string
	SERVER_PROTOCOL string
	REQUEST_URI     string
}

var DefaultConfig Config

func init() {
	DefaultConfig = Config{
		Debug:      false,
		ConfigFile: "etc/config.yaml",
		Serve:      ServeConfig{Location: "localhost:3000"},
		ServeCGI: ServeCGIConfig{
			REQUEST_METHOD:  http.MethodGet,
			SERVER_PROTOCOL: "HTTP/1.1",
			REQUEST_URI:     "/",
		},
	}
}
