// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.21.12
// source: message/message.proto

package message

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

// MessageQueueClient is the client API for MessageQueue service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type MessageQueueClient interface {
	Send(ctx context.Context, in *SendMsg, opts ...grpc.CallOption) (*ReplyMsg, error)
}

type messageQueueClient struct {
	cc grpc.ClientConnInterface
}

func NewMessageQueueClient(cc grpc.ClientConnInterface) MessageQueueClient {
	return &messageQueueClient{cc}
}

func (c *messageQueueClient) Send(ctx context.Context, in *SendMsg, opts ...grpc.CallOption) (*ReplyMsg, error) {
	out := new(ReplyMsg)
	err := c.cc.Invoke(ctx, "/message.MessageQueue/Send", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// MessageQueueServer is the server API for MessageQueue service.
// All implementations must embed UnimplementedMessageQueueServer
// for forward compatibility
type MessageQueueServer interface {
	Send(context.Context, *SendMsg) (*ReplyMsg, error)
	mustEmbedUnimplementedMessageQueueServer()
}

// UnimplementedMessageQueueServer must be embedded to have forward compatible implementations.
type UnimplementedMessageQueueServer struct {
}

func (UnimplementedMessageQueueServer) Send(context.Context, *SendMsg) (*ReplyMsg, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Send not implemented")
}
func (UnimplementedMessageQueueServer) mustEmbedUnimplementedMessageQueueServer() {}

// UnsafeMessageQueueServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to MessageQueueServer will
// result in compilation errors.
type UnsafeMessageQueueServer interface {
	mustEmbedUnimplementedMessageQueueServer()
}

func RegisterMessageQueueServer(s grpc.ServiceRegistrar, srv MessageQueueServer) {
	s.RegisterService(&MessageQueue_ServiceDesc, srv)
}

func _MessageQueue_Send_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SendMsg)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MessageQueueServer).Send(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/message.MessageQueue/Send",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MessageQueueServer).Send(ctx, req.(*SendMsg))
	}
	return interceptor(ctx, in, info, handler)
}

// MessageQueue_ServiceDesc is the grpc.ServiceDesc for MessageQueue service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var MessageQueue_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "message.MessageQueue",
	HandlerType: (*MessageQueueServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Send",
			Handler:    _MessageQueue_Send_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "message/message.proto",
}