package server

import (
	"net"
	"strconv"
)

type Conf struct {
	Host string
	Port int
}

func (s Conf) GetAddress() string {
	return net.JoinHostPort(s.Host, strconv.Itoa(s.Port))
}

type MetricsConf struct {
	LoadAverage  bool `toml:"load_average"`
	CPULoad      bool `toml:"cpu_load"`
	DiskLoad     bool `toml:"disk_load"`
	DiskInfo     bool `toml:"disk_info"`
	NetworkStats bool `toml:"network_stats"`
}
