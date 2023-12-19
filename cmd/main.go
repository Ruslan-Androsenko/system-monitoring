package main

import (
	"flag"

	"github.com/Ruslan-Androsenko/system-monitoring/logger"
)

var (
	configFile string
	logg       *logger.Logger
)

func init() {
	flag.StringVar(&configFile, "config", "/etc/system-monitoring/config.toml", "Path to configuration file")
}

func main() {
	flag.Parse()

	if hasVersionCommand() {
		printVersion()
		return
	}

	config := NewConfig()
	logg = logger.New(config.Logger.Level)
}
