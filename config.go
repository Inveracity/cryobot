package main

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
)

/*
this site is awesome: https://transform.now.sh/json-to-go/
*/

// Config contains the variables required by the config.yml file
type Config struct {
	Twitter struct {
		Track     []string `yaml:"track"`
		Follow    []string `yaml:"follow"`
		Blacklist []string `yaml:"blacklist"`
	} `yaml:"twitter"`
}

// ReadConfig reads the config.yml file and returns a Config struct
func ReadConfig() Config {
	b, err := ioutil.ReadFile("./config.yml")
	if err != nil {
		log.Println("Cannot find config.yml")
		log.Println("copy config_example.yml to config.yml and try again")
		os.Exit(1)
	}

	var cfg Config
	if err := yaml.Unmarshal(b, &cfg); err != nil {
		log.Fatal(err)
	}
	log.Println("Loading config config.yml")
	log.Println(cfg)

	return cfg
}
