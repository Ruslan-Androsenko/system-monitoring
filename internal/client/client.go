package client

import (
	"context"
	"time"

	"github.com/Ruslan-Androsenko/system-monitoring/api/proto"
	"github.com/Ruslan-Androsenko/system-monitoring/internal/logger"
	"github.com/Ruslan-Androsenko/system-monitoring/internal/server"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func InitGrpcClient(ctx context.Context, config server.Conf, logg *logger.Logger) {
	conn, err := grpc.Dial(config.GetAddress(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logg.Fatalf("Can not open connection: %v", err)
	}

	defer func() {
		err = conn.Close()
		if err != nil {
			logg.Fatalf("Can not close connection: %v", err)
		}
	}()

	client := proto.NewSystemMonitoringClient(conn)
	stream, err := client.Metrics(ctx, &proto.MonitoringRequest{
		EverySeconds: 5,
		AvgSeconds:   15,
	})
	if err != nil {
		logg.Fatalf("Can not creating stream: %v", err)
	}

	for i := 0; i < 100; i++ {
		select {
		case <-ctx.Done():
			return

		default:
			response, err := stream.Recv()
			if err != nil {
				logg.Fatalf("Can not receiving: %v", err)
			}

			logg.Infof("response: %v", response)
			time.Sleep(1 * time.Second)
		}
	}
}
