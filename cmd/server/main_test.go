package main_test

import (
	"context"
	"os/signal"
	"syscall"
	"testing"

	"github.com/Ruslan-Androsenko/system-monitoring/api/proto"
	"github.com/Ruslan-Androsenko/system-monitoring/internal/server"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	avgSeconds    = 15
	everySeconds  = 5
	countMessages = 10

	serverHost = "localhost"
	serverPort = 8070
	zeroNumber = 0.0
)

var metricsConf server.MetricsConf

func init() {
	metricsConf.LoadAverage = true
	metricsConf.CPULoad = true
	metricsConf.DiskLoad = true
	metricsConf.DiskInfo = true
	metricsConf.NetworkStats = false
}

func TestIntegration(t *testing.T) {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	config := server.Conf{Host: serverHost, Port: serverPort}
	conn, err := grpc.Dial(config.GetAddress(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(t, err)

	defer func() {
		err = conn.Close()
		require.NoError(t, err)
	}()

	client := proto.NewSystemMonitoringClient(conn)
	stream, err := client.Metrics(ctx, &proto.MonitoringRequest{
		EverySeconds: everySeconds,
		AvgSeconds:   avgSeconds,
	})
	require.NoError(t, err)

	var countReceived int

	for i := 0; i < countMessages; i++ {
		select {
		case <-ctx.Done():
			break

		default:
			response, err := stream.Recv()
			require.NoError(t, err)

			dataItemTest(t, response)
			countReceived++
		}
	}

	require.Equal(t, countMessages, countReceived)
}

// Проверка правильности заполнения данных по метрикам.
func dataItemTest(t *testing.T, dataItem *proto.MonitoringResponse) {
	t.Helper()

	if metricsConf.LoadAverage {
		require.GreaterOrEqual(t, dataItem.LoadAverage, zeroNumber)
	}

	// Проверяем заполненность данными для загрузки процессора
	if metricsConf.CPULoad {
		require.NotNil(t, dataItem.CpuLoad)
		require.GreaterOrEqual(t, dataItem.CpuLoad.UserMode, zeroNumber)
		require.GreaterOrEqual(t, dataItem.CpuLoad.SystemMode, zeroNumber)
		require.GreaterOrEqual(t, dataItem.CpuLoad.Idle, zeroNumber)
	}

	// Проверяем заполненность данными для загрузки диска
	if metricsConf.DiskLoad {
		require.NotNil(t, dataItem.DiskLoad)
		require.GreaterOrEqual(t, dataItem.DiskLoad.TransferPerSecond, zeroNumber)
		require.GreaterOrEqual(t, dataItem.DiskLoad.ReadPerSecond, zeroNumber)
		require.GreaterOrEqual(t, dataItem.DiskLoad.WritePerSecond, zeroNumber)
	}

	// Проверяем заполненность данными для информации об использовании диска
	if metricsConf.DiskInfo {
		require.NotNil(t, dataItem.DiskInfo)
		require.GreaterOrEqual(t, len(dataItem.DiskInfo), 0)
	}

	// Проверяем заполненность данными для сетевой статистики
	if metricsConf.NetworkStats {
		require.NotNil(t, dataItem.NetworkStats)
		require.NotNil(t, dataItem.NetworkStats.ListenerSocket)
		require.GreaterOrEqual(t, len(dataItem.NetworkStats.ListenerSocket), 0)

		// Проверяем заполненность данными для количества соединений
		require.NotNil(t, dataItem.NetworkStats.CounterConnections)
		require.NotNil(t, dataItem.NetworkStats.CounterConnections.Tcp)
		require.NotNil(t, dataItem.NetworkStats.CounterConnections.Udp)

		// Проверяем заполненность данными для
		require.GreaterOrEqual(t, len(dataItem.NetworkStats.CounterConnections.Tcp), 0)
		require.GreaterOrEqual(t, len(dataItem.NetworkStats.CounterConnections.Udp), 0)
	}
}
