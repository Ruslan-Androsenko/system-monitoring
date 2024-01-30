package config

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

type SetupConf struct {
	PathFile   string
	ServerHost string
	ServerPort int
}

func NewConfig(conf SetupConf) Config {
	var config Config

	if _, err := toml.DecodeFile(conf.PathFile, &config); err != nil {
		log.Fatalf("Can not read config file, err: %v \n", err)
	}

	if config.Server.Override {
		// Если через флаг передан хост, то используем его
		if config.Server.Host != conf.ServerHost {
			config.Server.Host = conf.ServerHost
		}

		// Если через флаг передан порт, то используем его
		if config.Server.Port != conf.ServerPort {
			config.Server.Port = conf.ServerPort
		}
	}

	return config
}
