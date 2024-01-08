package server

import "github.com/Ruslan-Androsenko/system-monitoring/api/proto"

type MetricsChannels struct {
	loadAverageCh chan float64
	cpuLoadCh     chan *proto.CpuLoad
	diskLoadCh    chan *proto.DiskLoad
	diskInfoCh    chan map[string]*proto.DiskInfo
}

// Инициализируем каналы для получения метрик системы.
func (m *MetricsChannels) init() {
	m.loadAverageCh = make(chan float64)
	m.cpuLoadCh = make(chan *proto.CpuLoad)
	m.diskLoadCh = make(chan *proto.DiskLoad)
	m.diskInfoCh = make(chan map[string]*proto.DiskInfo)
}
