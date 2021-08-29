// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package export1

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

// StoreClient is the client API for Store service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type StoreClient interface {
	/// Series streams each Series (Labels and chunk/downsampling chunk) for given label matchers and time range.
	///
	/// Series should strictly stream full series after series, optionally split by time. This means that a single frame can contain
	/// partition of the single series, but once a new series is started to be streamed it means that no more data will
	/// be sent for previous one.
	/// Series has to be sorted.
	///
	/// There is no requirements on chunk sorting, however it is recommended to have chunk sorted by chunk min time.
	/// This heavily optimizes the resource usage on Querier / Federated Queries.
	Series(ctx context.Context, in *SeriesRequest, opts ...grpc.CallOption) (Store_SeriesClient, error)
}

type storeClient struct {
	cc grpc.ClientConnInterface
}

func NewStoreClient(cc grpc.ClientConnInterface) StoreClient {
	return &storeClient{cc}
}

func (c *storeClient) Series(ctx context.Context, in *SeriesRequest, opts ...grpc.CallOption) (Store_SeriesClient, error) {
	stream, err := c.cc.NewStream(ctx, &Store_ServiceDesc.Streams[0], "/thanos.Store/Series", opts...)
	if err != nil {
		return nil, err
	}
	x := &storeSeriesClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type Store_SeriesClient interface {
	Recv() (*SeriesResponse, error)
	grpc.ClientStream
}

type storeSeriesClient struct {
	grpc.ClientStream
}

func (x *storeSeriesClient) Recv() (*SeriesResponse, error) {
	m := new(SeriesResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// StoreServer is the server API for Store service.
// All implementations should embed UnimplementedStoreServer
// for forward compatibility
type StoreServer interface {
	/// Series streams each Series (Labels and chunk/downsampling chunk) for given label matchers and time range.
	///
	/// Series should strictly stream full series after series, optionally split by time. This means that a single frame can contain
	/// partition of the single series, but once a new series is started to be streamed it means that no more data will
	/// be sent for previous one.
	/// Series has to be sorted.
	///
	/// There is no requirements on chunk sorting, however it is recommended to have chunk sorted by chunk min time.
	/// This heavily optimizes the resource usage on Querier / Federated Queries.
	Series(*SeriesRequest, Store_SeriesServer) error
}

// UnimplementedStoreServer should be embedded to have forward compatible implementations.
type UnimplementedStoreServer struct {
}

func (UnimplementedStoreServer) Series(*SeriesRequest, Store_SeriesServer) error {
	return status.Errorf(codes.Unimplemented, "method Series not implemented")
}

// UnsafeStoreServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to StoreServer will
// result in compilation errors.
type UnsafeStoreServer interface {
	mustEmbedUnimplementedStoreServer()
}

func RegisterStoreServer(s grpc.ServiceRegistrar, srv StoreServer) {
	s.RegisterService(&Store_ServiceDesc, srv)
}

func _Store_Series_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(SeriesRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(StoreServer).Series(m, &storeSeriesServer{stream})
}

type Store_SeriesServer interface {
	Send(*SeriesResponse) error
	grpc.ServerStream
}

type storeSeriesServer struct {
	grpc.ServerStream
}

func (x *storeSeriesServer) Send(m *SeriesResponse) error {
	return x.ServerStream.SendMsg(m)
}

// Store_ServiceDesc is the grpc.ServiceDesc for Store service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Store_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "thanos.Store",
	HandlerType: (*StoreServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Series",
			Handler:       _Store_Series_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "rpc.proto",
}
