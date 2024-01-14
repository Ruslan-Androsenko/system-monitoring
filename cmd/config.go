package main

import (
	"log"

	"github.com/BurntSushi/toml"
	"github.com/Ruslan-Androsenko/system-monitoring/internal/server"
)

type Config struct {
	Logger  LoggerConf
	Metrics server.MetricsConf
	Server  server.Conf
}

type LoggerConf struct {
	Level string
}

func NewConfig() Config {
	var config Config

	if _, err := toml.DecodeFile(configFile, &config); err != nil {
		log.Fatalf("Can not read config file, err: %v \n", err)
	}

	if config.Server.Override {
		// Если через флаг передан хост, то используем его
		if config.Server.Host != serverHost {
			config.Server.Host = serverHost
		}

		// Если через флаг передан порт, то используем его
		if config.Server.Port != serverPort {
			config.Server.Port = serverPort
		}
	}

	return config
}
