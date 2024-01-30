package constants

import (
	"os"
)

var (
	GithubKey  string
	APIPort    string
	CorsDomain string
)

func init() {
	GithubKey = os.Getenv("GITHUB_KEY")
	if GithubKey == "" {
		panic("GITHUB_KEY was not set")
	}
	APIPort = os.Getenv("API_PORT")
	if APIPort == "" {
		panic("API_PORT was not set")
	}

	CorsDomain = os.Getenv("CORS_DOMAIN")
	if CorsDomain == "" {
		panic("CORS_DOMAIN was not set")
	}
}
