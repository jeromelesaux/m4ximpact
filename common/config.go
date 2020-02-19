package common

import (
	"encoding/json"
	"fmt"
	"os"
)

var configFilepath = "m4config.json"
var defaultConfig = &Config{M4Url: "cpc"}

type Config struct {
	M4Url     string
	MailerApp string
}

func NewConfig() *Config {

	s, err := os.Stat(configFilepath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error while reading configuration file %s error %v\n", configFilepath, err)
		return defaultConfig
	}
	if !s.IsDir() {
		f, err := os.Open(configFilepath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error while opening configuration file %s error %v\n", configFilepath, err)
			return defaultConfig
		}
		defer f.Close()
		conf := defaultConfig
		if err := json.NewDecoder(f).Decode(conf); err != nil {
			fmt.Fprintf(os.Stderr, "Error while decoding configuration file %s error %v\n", configFilepath, err)
			return conf
		}
		return conf
	}
	return defaultConfig
}

func (c *Config) Save() error {
	f, err := os.Create(configFilepath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error while creating configuration file %s error %v\n", configFilepath, err)
		return err
	}
	defer f.Close()
	return json.NewEncoder(f).Encode(c)
}
