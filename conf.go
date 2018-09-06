package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"path/filepath"
)

// Conf is for deserializing the JSON config file.
type Conf struct {
	Token  string `json:"token"`
	Prefix string `json:"prefix"`
	Status string `json:"string"`
	State  map[string]interface{}
}

// LoadConfig instantiates a Viper object with config info required for the bot to work.
func LoadConfig() Conf {
	var v Conf

	user, err := user.Current()
	if err != nil {
		log.Fatalf("Unable to load config file: %s", err)
	}

	home := user.HomeDir

	f, err := os.Open(filepath.Join(home, "alpha.json"))
	if err != nil {
		log.Fatalf("Unable to load config file: %s", err)
	}

	rawConf, err := ioutil.ReadAll(f)
	if err != nil {
		log.Fatalf("Unable to read config file: %s", err)
	}

	json.Unmarshal(rawConf, &v)

	return v
}
