// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v3.12.4
// source: msgpb/message.proto

package msgpb

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
	CommunicationService_SendMessage_FullMethodName = "/msg.CommunicationService/SendMessage"
	CommunicationService_GetMessages_FullMethodName = "/msg.CommunicationService/GetMessages"
)

// CommunicationServiceClient is the client API for CommunicationService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type CommunicationServiceClient interface {
	SendMessage(ctx context.Context, in *Message, opts ...grpc.CallOption) (*Message, error)
	GetMessages(ctx context.Context, in *NoParams, opts ...grpc.CallOption) (CommunicationService_GetMessagesClient, error)
}

type communicationServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewCommunicationServiceClient(cc grpc.ClientConnInterface) CommunicationServiceClient {
	return &communicationServiceClient{cc}
}

func (c *communicationServiceClient) SendMessage(ctx context.Context, in *Message, opts ...grpc.CallOption) (*Message, error) {
	out := new(Message)
	err := c.cc.Invoke(ctx, CommunicationService_SendMessage_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *communicationServiceClient) GetMessages(ctx context.Context, in *NoParams, opts ...grpc.CallOption) (CommunicationService_GetMessagesClient, error) {
	stream, err := c.cc.NewStream(ctx, &CommunicationService_ServiceDesc.Streams[0], CommunicationService_GetMessages_FullMethodName, opts...)
	if err != nil {
		return nil, err
	}
	x := &communicationServiceGetMessagesClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type CommunicationService_GetMessagesClient interface {
	Recv() (*Message, error)
	grpc.ClientStream
}

type communicationServiceGetMessagesClient struct {
	grpc.ClientStream
}

func (x *communicationServiceGetMessagesClient) Recv() (*Message, error) {
	m := new(Message)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// CommunicationServiceServer is the server API for CommunicationService service.
// All implementations must embed UnimplementedCommunicationServiceServer
// for forward compatibility
type CommunicationServiceServer interface {
	SendMessage(context.Context, *Message) (*Message, error)
	GetMessages(*NoParams, CommunicationService_GetMessagesServer) error
	mustEmbedUnimplementedCommunicationServiceServer()
}

// UnimplementedCommunicationServiceServer must be embedded to have forward compatible implementations.
type UnimplementedCommunicationServiceServer struct {
}

func (UnimplementedCommunicationServiceServer) SendMessage(context.Context, *Message) (*Message, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SendMessage not implemented")
}
func (UnimplementedCommunicationServiceServer) GetMessages(*NoParams, CommunicationService_GetMessagesServer) error {
	return status.Errorf(codes.Unimplemented, "method GetMessages not implemented")
}
func (UnimplementedCommunicationServiceServer) mustEmbedUnimplementedCommunicationServiceServer() {}

// UnsafeCommunicationServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to CommunicationServiceServer will
// result in compilation errors.
type UnsafeCommunicationServiceServer interface {
	mustEmbedUnimplementedCommunicationServiceServer()
}

func RegisterCommunicationServiceServer(s grpc.ServiceRegistrar, srv CommunicationServiceServer) {
	s.RegisterService(&CommunicationService_ServiceDesc, srv)
}

func _CommunicationService_SendMessage_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Message)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CommunicationServiceServer).SendMessage(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: CommunicationService_SendMessage_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CommunicationServiceServer).SendMessage(ctx, req.(*Message))
	}
	return interceptor(ctx, in, info, handler)
}

func _CommunicationService_GetMessages_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(NoParams)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(CommunicationServiceServer).GetMessages(m, &communicationServiceGetMessagesServer{stream})
}

type CommunicationService_GetMessagesServer interface {
	Send(*Message) error
	grpc.ServerStream
}

type communicationServiceGetMessagesServer struct {
	grpc.ServerStream
}

func (x *communicationServiceGetMessagesServer) Send(m *Message) error {
	return x.ServerStream.SendMsg(m)
}

// CommunicationService_ServiceDesc is the grpc.ServiceDesc for CommunicationService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var CommunicationService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "msg.CommunicationService",
	HandlerType: (*CommunicationServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SendMessage",
			Handler:    _CommunicationService_SendMessage_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "GetMessages",
			Handler:       _CommunicationService_GetMessages_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "msgpb/message.proto",
}
