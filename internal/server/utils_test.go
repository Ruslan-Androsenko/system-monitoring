package server

import (
	"context"
	"testing"
	"time"

	"github.com/Ruslan-Androsenko/system-monitoring/api/proto"
	"github.com/stretchr/testify/require"
)

const (
	avgSeconds   = 15
	everySeconds = 5
	zeroNumber   = 0.0
)

func init() {
	metricsConf.LoadAverage = true
	metricsConf.CPULoad = true
	metricsConf.DiskLoad = true
	metricsConf.DiskInfo = true
	metricsConf.NetworkStats = false
}

// Проверка инициализации каналов.
func initMetricsChannelsTest(t *testing.T, metricsChs MetricsChannels) {
	t.Helper()

	require.NotNil(t, metricsChs.errCh)
	require.NotNil(t, metricsChs.loadAverageCh)
	require.NotNil(t, metricsChs.cpuLoadCh)
	require.NotNil(t, metricsChs.diskLoadCh)
	require.NotNil(t, metricsChs.diskInfoCh)
	require.NotNil(t, metricsChs.networkStatsCh)
}

// Проверка инициализации мьютексов.
func initMetricsMutexesTest(t *testing.T, mu MetricsMutex) {
	t.Helper()

	require.NotNil(t, mu.loadAverage)
	require.NotNil(t, mu.cpuLoad)
	require.NotNil(t, mu.diskLoad)
	require.NotNil(t, mu.diskInfo)
	require.NotNil(t, mu.networkStats)
}

// Проверка правильности заполнения данных по метрикам.
func dataItemTest(t *testing.T, mu *MetricsMutex, dataItem *proto.MonitoringResponse) {
	t.Helper()

	if metricsConf.LoadAverage {
		mu.loadAverage.RLock()
		require.GreaterOrEqual(t, dataItem.LoadAverage, zeroNumber)
		mu.loadAverage.RUnlock()
	}

	// Проверяем заполненность данными для загрузки процессора
	if metricsConf.CPULoad {
		mu.cpuLoad.RLock()
		require.NotNil(t, dataItem.CpuLoad)
		require.GreaterOrEqual(t, dataItem.CpuLoad.UserMode, zeroNumber)
		require.GreaterOrEqual(t, dataItem.CpuLoad.SystemMode, zeroNumber)
		require.GreaterOrEqual(t, dataItem.CpuLoad.Idle, zeroNumber)
		mu.cpuLoad.RUnlock()
	}

	// Проверяем заполненность данными для загрузки диска
	if metricsConf.DiskLoad {
		mu.diskLoad.RLock()
		require.NotNil(t, dataItem.DiskLoad)
		require.GreaterOrEqual(t, dataItem.DiskLoad.TransferPerSecond, zeroNumber)
		require.GreaterOrEqual(t, dataItem.DiskLoad.ReadPerSecond, zeroNumber)
		require.GreaterOrEqual(t, dataItem.DiskLoad.WritePerSecond, zeroNumber)
		mu.diskLoad.RUnlock()
	}

	// Проверяем заполненность данными для информации об использовании диска
	if metricsConf.DiskInfo {
		mu.diskInfo.RLock()
		require.NotNil(t, dataItem.DiskInfo)
		require.GreaterOrEqual(t, len(dataItem.DiskInfo), 0)
		mu.diskInfo.RUnlock()
	}

	// Проверяем заполненность данными для сетевой статистики
	if metricsConf.NetworkStats {
		mu.networkStats.RLock()
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
		mu.networkStats.RUnlock()
	}
}

func TestFillDataItem(t *testing.T) {
	var (
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
	dataItemTest(t, &mu, dataItem)
}

func TestMakeDataSlice(t *testing.T) {
	var (
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

	data := make([]*proto.MonitoringResponse, countSeconds)

	for i := 0; i < countSeconds; i++ {
		if data[i] == nil {
			data[i] = &proto.MonitoringResponse{}
		}

		data[i] = fillDataItem(data[i], metricsChs, mu)
		require.NotNil(t, data[i])

		if i >= avgSeconds && i%everySeconds == 0 {
			dataSlice := makeDataSlice(data, i, avgSeconds)
			require.NotNil(t, dataSlice)
			require.Equal(t, avgSeconds, len(dataSlice))

			response := calculateAverageOfSlice(dataSlice)
			require.NotNil(t, response)
			dataItemTest(t, &mu, response)
		}
	}
}
