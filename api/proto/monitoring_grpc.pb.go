// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v3.12.4
// source: proto/monitoring.proto

package proto

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

const (
	SystemMonitoring_Metrics_FullMethodName = "/SystemMonitoring/Metrics"
)

// SystemMonitoringClient is the client API for SystemMonitoring service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type SystemMonitoringClient interface {
	Metrics(ctx context.Context, in *MonitoringRequest, opts ...grpc.CallOption) (SystemMonitoring_MetricsClient, error)
}

type systemMonitoringClient struct {
	cc grpc.ClientConnInterface
}

func NewSystemMonitoringClient(cc grpc.ClientConnInterface) SystemMonitoringClient {
	return &systemMonitoringClient{cc}
}

func (c *systemMonitoringClient) Metrics(ctx context.Context, in *MonitoringRequest, opts ...grpc.CallOption) (SystemMonitoring_MetricsClient, error) {
	stream, err := c.cc.NewStream(ctx, &SystemMonitoring_ServiceDesc.Streams[0], SystemMonitoring_Metrics_FullMethodName, opts...)
	if err != nil {
		return nil, err
	}
	x := &systemMonitoringMetricsClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type SystemMonitoring_MetricsClient interface {
	Recv() (*MonitoringResponse, error)
	grpc.ClientStream
}

type systemMonitoringMetricsClient struct {
	grpc.ClientStream
}

func (x *systemMonitoringMetricsClient) Recv() (*MonitoringResponse, error) {
	m := new(MonitoringResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// SystemMonitoringServer is the server API for SystemMonitoring service.
// All implementations must embed UnimplementedSystemMonitoringServer
// for forward compatibility
type SystemMonitoringServer interface {
	Metrics(*MonitoringRequest, SystemMonitoring_MetricsServer) error
	mustEmbedUnimplementedSystemMonitoringServer()
}

// UnimplementedSystemMonitoringServer must be embedded to have forward compatible implementations.
type UnimplementedSystemMonitoringServer struct {
}

func (UnimplementedSystemMonitoringServer) Metrics(*MonitoringRequest, SystemMonitoring_MetricsServer) error {
	return status.Errorf(codes.Unimplemented, "method Metrics not implemented")
}
func (UnimplementedSystemMonitoringServer) mustEmbedUnimplementedSystemMonitoringServer() {}

// UnsafeSystemMonitoringServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to SystemMonitoringServer will
// result in compilation errors.
type UnsafeSystemMonitoringServer interface {
	mustEmbedUnimplementedSystemMonitoringServer()
}

func RegisterSystemMonitoringServer(s grpc.ServiceRegistrar, srv SystemMonitoringServer) {
	s.RegisterService(&SystemMonitoring_ServiceDesc, srv)
}

func _SystemMonitoring_Metrics_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(MonitoringRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(SystemMonitoringServer).Metrics(m, &systemMonitoringMetricsServer{stream})
}

type SystemMonitoring_MetricsServer interface {
	Send(*MonitoringResponse) error
	grpc.ServerStream
}

type systemMonitoringMetricsServer struct {
	grpc.ServerStream
}

func (x *systemMonitoringMetricsServer) Send(m *MonitoringResponse) error {
	return x.ServerStream.SendMsg(m)
}

// SystemMonitoring_ServiceDesc is the grpc.ServiceDesc for SystemMonitoring service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var SystemMonitoring_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "SystemMonitoring",
	HandlerType: (*SystemMonitoringServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Metrics",
			Handler:       _SystemMonitoring_Metrics_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "proto/monitoring.proto",
}
