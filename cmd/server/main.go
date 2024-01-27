package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/Ruslan-Androsenko/system-monitoring/internal/config"
	"github.com/Ruslan-Androsenko/system-monitoring/internal/logger"
	"github.com/Ruslan-Androsenko/system-monitoring/internal/server"
)

var (
	configFile string
	serverHost string
	serverPort int
	logg       *logger.Logger
)

func init() {
	flag.StringVar(&configFile, "config", "/etc/system-monitoring/config.toml", "Path to configuration file")
	flag.StringVar(&serverHost, "host", "localhost", "Host to start the server")
	flag.IntVar(&serverPort, "port", 8080, "Port to start the server")
}

func main() {
	flag.Parse()

	if hasVersionCommand() {
		printVersion()
		return
	}

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	appConfig := config.NewConfig(config.SetupConf{
		PathFile:   configFile,
		ServerHost: serverHost,
		ServerPort: serverPort,
	})
	logg = logger.New(appConfig.Logger.Level)
	grpcServer := server.NewServer(appConfig.Server, appConfig.Metrics, logg)

	go func() {
		<-ctx.Done()
		logg.Info("system-monitoring is stopped...")

		if err := grpcServer.Stop(); err != nil {
			logg.Error("failed to stop grpc server: " + err.Error())
		}
	}()

	logg.Info("system-monitoring is running...")

	if err := grpcServer.Start(); err != nil {
		logg.Error("failed to start grpc server: " + err.Error())
		cancel()
		os.Exit(1) //nolint:gocritic
	}
}
