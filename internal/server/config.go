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
