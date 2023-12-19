package main

import (
	"log"

	"github.com/BurntSushi/toml"
)

type Config struct {
	Logger LoggerConf
}

type LoggerConf struct {
	Level string
}

func NewConfig() Config {
	var config Config

	if _, err := toml.DecodeFile(configFile, &config); err != nil {
		log.Fatalf("Can not read config file, err: %v \n", err)
	}

	return config
}
