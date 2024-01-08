package server

import (
	"context"
	"net"
	"time"

	"github.com/Ruslan-Androsenko/system-monitoring/api/proto"
	"github.com/Ruslan-Androsenko/system-monitoring/internal/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const countSeconds = 60

var (
	logg        *logger.Logger
	metricsConf MetricsConf
)

type Server struct {
	proto.UnimplementedSystemMonitoringServer

	grpc     *grpc.Server
	listener net.Listener
}

func NewServer(server Conf, metrics MetricsConf, logger *logger.Logger) *Server {
	logg = logger
	metricsConf = metrics

	listener, err := net.Listen("tcp", server.GetAddress())
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
	var (
		cnt        int
		metricsChs MetricsChannels
	)

	data := make([]*proto.MonitoringResponse, countSeconds)
	avgSeconds := int(req.AvgSeconds)
	everySeconds := int(req.EverySeconds)
	ctx, cancel := context.WithCancel(context.Background())

	metricsChs.init()
	metricsChs.run(ctx)

	defer func() {
		cancel()
		close(metricsChs.errCh)
	}()

	go func() {
		err := <-metricsChs.errCh
		if err != nil {
			logg.Error(err.Error())
			cancel()
		}
	}()

	logg.Infof("request: %v", req)

	for i := 0; ; i++ {
		if i == countSeconds {
			i = 0
		}

		if data[i] == nil {
			data[i] = &proto.MonitoringResponse{}
		}

		data[i] = fillDataSlice(data[i], metricsChs)

		if cnt >= avgSeconds && cnt%everySeconds == 0 {
			dataSlice := makeDataSlice(data, i, avgSeconds)
			response := calculateAverageOfSlice(dataSlice)

			if err := stream.Send(response); err != nil {
				logg.Errorf("Can not send data to stream, response: %v. Error: %v", response, err)
				return err
			}
		}

		time.Sleep(1 * time.Second)
		cnt++
	}
}
