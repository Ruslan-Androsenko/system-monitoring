package main

import (
	"context"
	"flag"
	"os/signal"
	"syscall"

	"github.com/Ruslan-Androsenko/system-monitoring/internal/config"
	"github.com/Ruslan-Androsenko/system-monitoring/internal/logger"
)

var (
	configFile   string
	serverHost   string
	serverPort   int
	messages     int
	everySeconds int
	avgSeconds   int
	logg         *logger.Logger
)

func init() {
	flag.StringVar(&configFile, "config", "/etc/system-monitoring/config.toml", "Path to configuration file")
	flag.StringVar(&serverHost, "host", "localhost", "Host to connect to the server")
	flag.IntVar(&serverPort, "port", 8080, "Port to connect to the server")
	flag.IntVar(&messages, "messages", 50, "Number of messages received")
	flag.IntVar(&everySeconds, "everySeconds", 5, "Every second")
	flag.IntVar(&avgSeconds, "avgSeconds", 15, "Average second")
}

func main() {
	flag.Parse()

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	appConfig := config.NewConfig(config.SetupConf{
		PathFile:   configFile,
		ServerHost: serverHost,
		ServerPort: serverPort,
	})
	logg = logger.New(appConfig.Logger.Level)

	go func() {
		<-ctx.Done()
		logg.Info("Grpc client is stopped...")
	}()

	logg.Info("Grpc client is receiving...")
	initGrpcClient(ctx, appConfig.Server)
}
