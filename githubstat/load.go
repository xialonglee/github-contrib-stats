package githubstat

import (
	"encoding/json"
	"io/ioutil"
)

type Load interface {
}

var config Config

type Config struct {
	AccessToken string
	OrgName     string
}

// read config file
var _ = func() int {
	file, err := ioutil.ReadFile("./config.json")
	if err != nil {
		panic(err)
	}
	json.Unmarshal(file, &config)
	return 0
}()
