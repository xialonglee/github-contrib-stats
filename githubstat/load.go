package githubstat

import (
	"time"

	"github.com/BurntSushi/toml"
)

type Load interface {
}

var Config Configuration

type Configuration struct {
	StatBeginTime time.Time
	AccessToken   string
	Users         []string
	Repos         []string
	Metrics       string
}

// read config file
var _ = func() int {
	if _, err := toml.DecodeFile("./config.toml", &Config); err != nil {
		panic(err)
	}
	return 0
}()
