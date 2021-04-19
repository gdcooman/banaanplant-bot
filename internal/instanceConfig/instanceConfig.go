package instanceConfig

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Token string `yaml:"token"`
	Dsn   string `yaml:"dsn"`
}

func ParseConfigFromYAMLFile(file string) (c *Config, err error) {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatal(err)
	}

	err = yaml.Unmarshal(data, &c)
	return
}
