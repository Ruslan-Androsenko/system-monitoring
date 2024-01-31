package main

import (
	"context"

	"github.com/Ruslan-Androsenko/system-monitoring/api/proto"
	"github.com/Ruslan-Androsenko/system-monitoring/internal/server"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func initGrpcClient(ctx context.Context, config server.Conf) {
	conn, err := grpc.Dial(config.GetAddress(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logg.Fatalf("Can not open connection: %v", err)
	}

	defer func() {
		if err = conn.Close(); err != nil {
			logg.Fatalf("Can not close connection: %v", err)
		}
	}()

	client := proto.NewSystemMonitoringClient(conn)
	stream, err := client.Metrics(ctx, &proto.MonitoringRequest{
		EverySeconds: uint32(everySeconds),
		AvgSeconds:   uint32(avgSeconds),
	})
	if err != nil {
		logg.Fatalf("Can not creating stream: %v", err)
	}

	for i := 0; i < messages; i++ {
		select {
		case <-ctx.Done():
			return

		default:
			response, err := stream.Recv()
			if err != nil {
				logg.Fatalf("Can not receiving: %v", err)
			}

			logg.Infof("response: %v \n", response)
		}
	}
}
