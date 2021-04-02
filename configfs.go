package goblog

import (
	"embed"

	"gopkg.in/yaml.v3"
)

var Conf Config

func init() {
	var conf Config
	f, err := ConfigFS.Open("config/config.yaml")
	if err != nil {
		panic("could not open config: " + err.Error())
	}
	defer f.Close()

	err = yaml.NewDecoder(f).Decode(&conf)
	if err != nil {
		panic("could not yaml decode config: " + err.Error())
	}
	Conf = conf
}

//go:embed config
var ConfigFS embed.FS
