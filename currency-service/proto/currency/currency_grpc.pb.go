// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.21.12
// source: currency.proto

package currency

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

// CurrencyClient is the client API for Currency service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type CurrencyClient interface {
	GetCurrencyRate(ctx context.Context, in *GetCurrencyRateRequest, opts ...grpc.CallOption) (*GetCurrencyRateResponse, error)
	StreamCurrencyRates(ctx context.Context, opts ...grpc.CallOption) (Currency_StreamCurrencyRatesClient, error)
}

type currencyClient struct {
	cc grpc.ClientConnInterface
}

func NewCurrencyClient(cc grpc.ClientConnInterface) CurrencyClient {
	return &currencyClient{cc}
}

func (c *currencyClient) GetCurrencyRate(ctx context.Context, in *GetCurrencyRateRequest, opts ...grpc.CallOption) (*GetCurrencyRateResponse, error) {
	out := new(GetCurrencyRateResponse)
	err := c.cc.Invoke(ctx, "/currency.Currency/GetCurrencyRate", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *currencyClient) StreamCurrencyRates(ctx context.Context, opts ...grpc.CallOption) (Currency_StreamCurrencyRatesClient, error) {
	stream, err := c.cc.NewStream(ctx, &Currency_ServiceDesc.Streams[0], "/currency.Currency/StreamCurrencyRates", opts...)
	if err != nil {
		return nil, err
	}
	x := &currencyStreamCurrencyRatesClient{stream}
	return x, nil
}

type Currency_StreamCurrencyRatesClient interface {
	Send(*StreamCurrencyRateRequest) error
	Recv() (*StreamCurrencyRateResponse, error)
	grpc.ClientStream
}

type currencyStreamCurrencyRatesClient struct {
	grpc.ClientStream
}

func (x *currencyStreamCurrencyRatesClient) Send(m *StreamCurrencyRateRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *currencyStreamCurrencyRatesClient) Recv() (*StreamCurrencyRateResponse, error) {
	m := new(StreamCurrencyRateResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// CurrencyServer is the server API for Currency service.
// All implementations must embed UnimplementedCurrencyServer
// for forward compatibility
type CurrencyServer interface {
	GetCurrencyRate(context.Context, *GetCurrencyRateRequest) (*GetCurrencyRateResponse, error)
	StreamCurrencyRates(Currency_StreamCurrencyRatesServer) error
	mustEmbedUnimplementedCurrencyServer()
}

// UnimplementedCurrencyServer must be embedded to have forward compatible implementations.
type UnimplementedCurrencyServer struct {
}

func (UnimplementedCurrencyServer) GetCurrencyRate(context.Context, *GetCurrencyRateRequest) (*GetCurrencyRateResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetCurrencyRate not implemented")
}
func (UnimplementedCurrencyServer) StreamCurrencyRates(Currency_StreamCurrencyRatesServer) error {
	return status.Errorf(codes.Unimplemented, "method StreamCurrencyRates not implemented")
}
func (UnimplementedCurrencyServer) mustEmbedUnimplementedCurrencyServer() {}

// UnsafeCurrencyServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to CurrencyServer will
// result in compilation errors.
type UnsafeCurrencyServer interface {
	mustEmbedUnimplementedCurrencyServer()
}

func RegisterCurrencyServer(s grpc.ServiceRegistrar, srv CurrencyServer) {
	s.RegisterService(&Currency_ServiceDesc, srv)
}

func _Currency_GetCurrencyRate_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetCurrencyRateRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CurrencyServer).GetCurrencyRate(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/currency.Currency/GetCurrencyRate",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CurrencyServer).GetCurrencyRate(ctx, req.(*GetCurrencyRateRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Currency_StreamCurrencyRates_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(CurrencyServer).StreamCurrencyRates(&currencyStreamCurrencyRatesServer{stream})
}

type Currency_StreamCurrencyRatesServer interface {
	Send(*StreamCurrencyRateResponse) error
	Recv() (*StreamCurrencyRateRequest, error)
	grpc.ServerStream
}

type currencyStreamCurrencyRatesServer struct {
	grpc.ServerStream
}

func (x *currencyStreamCurrencyRatesServer) Send(m *StreamCurrencyRateResponse) error {
	return x.ServerStream.SendMsg(m)
}

func (x *currencyStreamCurrencyRatesServer) Recv() (*StreamCurrencyRateRequest, error) {
	m := new(StreamCurrencyRateRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// Currency_ServiceDesc is the grpc.ServiceDesc for Currency service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Currency_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "currency.Currency",
	HandlerType: (*CurrencyServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetCurrencyRate",
			Handler:    _Currency_GetCurrencyRate_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "StreamCurrencyRates",
			Handler:       _Currency_StreamCurrencyRates_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "currency.proto",
}
