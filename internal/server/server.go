package server

import (
	"net"
	"time"

	"github.com/Ruslan-Androsenko/system-monitoring/api/proto"
	"github.com/Ruslan-Androsenko/system-monitoring/internal/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type Server struct {
	proto.UnimplementedSystemMonitoringServer

	grpc     *grpc.Server
	listener net.Listener
}

var logg *logger.Logger

func NewServer(config Conf, logger *logger.Logger) *Server {
	logg = logger

	listener, err := net.Listen("tcp", config.GetAddress())
	if err != nil {
		logg.Fatalf("Failed to listen: %v", err)
	}

	return &Server{
		grpc:     grpc.NewServer(),
		listener: listener,
	}
}

func (s *Server) Start() error {
	// Register service on gRPC server
	proto.RegisterSystemMonitoringServer(s.grpc, s)

	// Register reflection service on gRPC server
	reflection.Register(s.grpc)

	return s.grpc.Serve(s.listener)
}

func (s *Server) Stop() error {
	s.grpc.Stop()
	return s.listener.Close()
}

func (s *Server) Metrics(req *proto.MonitoringRequest, stream proto.SystemMonitoring_MetricsServer) error {
	cnt := 1
	logg.Infof("request: %v", req)

	for cnt < 100 {
		cnt++
		response := &proto.MonitoringResponse{
			LoadAverage: 123.76,
			CpuLoad: &proto.CpuLoad{
				UserMode:   111.1,
				SystemMode: 222.2,
				Idle:       333.3,
			},
			DiskLoad: &proto.DiscLoad{
				TransferPerSecond: 444.1,
				KbsPerSecond:      555.2,
			},
			DiskInfo: &proto.DiscInfo{
				UsageSize:  777.1,
				UsageInode: 888.2,
			},
		}

		if err := stream.Send(response); err != nil {
			return err
		}

		logg.Infof("response: %v", response)
		time.Sleep(1 * time.Second)
	}

	return nil
}
