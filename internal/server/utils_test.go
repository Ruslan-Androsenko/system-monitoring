package server

import (
	"context"
	"testing"
	"time"

	"github.com/Ruslan-Androsenko/system-monitoring/api/proto"
	"github.com/stretchr/testify/require"
)

func init() {
	metricsConf.LoadAverage = true
	metricsConf.CPULoad = true
	metricsConf.DiskLoad = true
	metricsConf.DiskInfo = true
	metricsConf.NetworkStats = true
}

func initMetricsChannelsTest(t *testing.T, metricsChs MetricsChannels) {
	t.Helper()
	require.NotNil(t, metricsChs.errCh)
	require.NotNil(t, metricsChs.loadAverageCh)
	require.NotNil(t, metricsChs.cpuLoadCh)
	require.NotNil(t, metricsChs.diskLoadCh)
	require.NotNil(t, metricsChs.diskInfoCh)
	require.NotNil(t, metricsChs.networkStatsCh)
}

func initMetricsMutexesTest(t *testing.T, mu MetricsMutex) {
	t.Helper()
	require.NotNil(t, mu.loadAverage)
	require.NotNil(t, mu.cpuLoad)
	require.NotNil(t, mu.diskLoad)
	require.NotNil(t, mu.diskInfo)
	require.NotNil(t, mu.networkStats)
}

func TestFillDataItem(t *testing.T) {
	var (
		zeroNumber float64
		metricsChs MetricsChannels
		mu         MetricsMutex
	)

	metricsChs.init()
	initMetricsChannelsTest(t, metricsChs)

	mu.init()
	initMetricsMutexesTest(t, mu)

	ctx, cancel := context.WithCancel(context.Background())
	metricsChs.run(ctx)

	defer func() {
		cancel()
		close(metricsChs.errCh)
	}()

	dataItem := fillDataItem(&proto.MonitoringResponse{}, metricsChs, mu)
	time.Sleep(time.Second)
	require.NotNil(t, dataItem)

	mu.loadAverage.RLock()
	require.Greater(t, dataItem.LoadAverage, zeroNumber)
	mu.loadAverage.RUnlock()

	// Проверяем заполненность данными для загрузки процессора
	mu.cpuLoad.RLock()
	require.NotNil(t, dataItem.CpuLoad)
	require.Greater(t, dataItem.CpuLoad.UserMode, zeroNumber)
	require.Greater(t, dataItem.CpuLoad.SystemMode, zeroNumber)
	require.Greater(t, dataItem.CpuLoad.Idle, zeroNumber)
	mu.cpuLoad.RUnlock()

	// Проверяем заполненность данными для загрузки диска
	mu.diskLoad.RLock()
	require.NotNil(t, dataItem.DiskLoad)
	require.Greater(t, dataItem.DiskLoad.TransferPerSecond, zeroNumber)
	require.Greater(t, dataItem.DiskLoad.ReadPerSecond, zeroNumber)
	require.Greater(t, dataItem.DiskLoad.WritePerSecond, zeroNumber)
	mu.diskLoad.RUnlock()

	// Проверяем заполненность данными для информации об использовании диска
	mu.diskInfo.RLock()
	require.NotNil(t, dataItem.DiskInfo)
	require.Greater(t, len(dataItem.DiskInfo), 0)
	mu.diskInfo.RUnlock()

	// Проверяем заполненность данными для сетевой статистики
	mu.networkStats.RLock()
	require.NotNil(t, dataItem.NetworkStats)
	require.NotNil(t, dataItem.NetworkStats.ListenerSocket)
	require.Greater(t, len(dataItem.NetworkStats.ListenerSocket), 0)

	// Проверяем заполненность данными для количества соединений
	require.NotNil(t, dataItem.NetworkStats.CounterConnections)
	require.NotNil(t, dataItem.NetworkStats.CounterConnections.Tcp)
	require.NotNil(t, dataItem.NetworkStats.CounterConnections.Udp)

	// Проверяем заполненность данными для
	require.Greater(t, len(dataItem.NetworkStats.CounterConnections.Tcp), 0)
	require.Greater(t, len(dataItem.NetworkStats.CounterConnections.Udp), 0)
	mu.networkStats.RUnlock()
}
