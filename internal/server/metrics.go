package server

import (
	"context"

	"github.com/Ruslan-Androsenko/system-monitoring/api/proto"
	"github.com/Ruslan-Androsenko/system-monitoring/internal/tools"
)

type MetricsChannels struct {
	errCh          chan error
	loadAverageCh  chan float64
	cpuLoadCh      chan *proto.CpuLoad
	diskLoadCh     chan *proto.DiskLoad
	diskInfoCh     chan map[string]*proto.DiskInfo
	networkStatsCh chan *proto.NetworkStats
}

// Инициализируем каналы для получения метрик системы.
func (m *MetricsChannels) init() {
	m.errCh = make(chan error)
	m.loadAverageCh = make(chan float64)
	m.cpuLoadCh = make(chan *proto.CpuLoad)
	m.diskLoadCh = make(chan *proto.DiskLoad)
	m.diskInfoCh = make(chan map[string]*proto.DiskInfo)
	m.networkStatsCh = make(chan *proto.NetworkStats)
}

// Запускаем сбор необходимых метрик системы.
func (m *MetricsChannels) run(ctx context.Context) {
	if metricsConf.LoadAverage {
		go tools.GetLoadAverage(ctx, m.loadAverageCh, m.errCh)
	}

	if metricsConf.CPULoad {
		go tools.GetCPULoad(ctx, m.cpuLoadCh, m.errCh)
	}

	if metricsConf.DiskLoad {
		go tools.GetDiskLoad(ctx, m.diskLoadCh, m.errCh)
	}

	if metricsConf.DiskInfo {
		go tools.GetDiskInfo(ctx, m.diskInfoCh, m.errCh)
	}

	if metricsConf.NetworkStats {
		go tools.GetNetworkStats(ctx, m.networkStatsCh, m.errCh)
	}
}
