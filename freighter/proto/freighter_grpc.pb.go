// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v3.14.0
// source: proto/freighter.proto

package freighter

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
	Freighter_GetDir_FullMethodName  = "/freighter.Freighter/GetDir"
	Freighter_GetFile_FullMethodName = "/freighter.Freighter/GetFile"
	Freighter_GetTree_FullMethodName = "/freighter.Freighter/GetTree"
)

// FreighterClient is the client API for Freighter service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type FreighterClient interface {
	GetDir(ctx context.Context, in *DirRequest, opts ...grpc.CallOption) (*DirReply, error)
	GetFile(ctx context.Context, in *FileRequest, opts ...grpc.CallOption) (*FileReply, error)
	GetTree(ctx context.Context, in *TreeRequest, opts ...grpc.CallOption) (*TreeReply, error)
}

type freighterClient struct {
	cc grpc.ClientConnInterface
}

func NewFreighterClient(cc grpc.ClientConnInterface) FreighterClient {
	return &freighterClient{cc}
}

func (c *freighterClient) GetDir(ctx context.Context, in *DirRequest, opts ...grpc.CallOption) (*DirReply, error) {
	out := new(DirReply)
	err := c.cc.Invoke(ctx, Freighter_GetDir_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *freighterClient) GetFile(ctx context.Context, in *FileRequest, opts ...grpc.CallOption) (*FileReply, error) {
	out := new(FileReply)
	err := c.cc.Invoke(ctx, Freighter_GetFile_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *freighterClient) GetTree(ctx context.Context, in *TreeRequest, opts ...grpc.CallOption) (*TreeReply, error) {
	out := new(TreeReply)
	err := c.cc.Invoke(ctx, Freighter_GetTree_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// FreighterServer is the server API for Freighter service.
// All implementations must embed UnimplementedFreighterServer
// for forward compatibility
type FreighterServer interface {
	GetDir(context.Context, *DirRequest) (*DirReply, error)
	GetFile(context.Context, *FileRequest) (*FileReply, error)
	GetTree(context.Context, *TreeRequest) (*TreeReply, error)
	mustEmbedUnimplementedFreighterServer()
}

// UnimplementedFreighterServer must be embedded to have forward compatible implementations.
type UnimplementedFreighterServer struct {
}

func (UnimplementedFreighterServer) GetDir(context.Context, *DirRequest) (*DirReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetDir not implemented")
}
func (UnimplementedFreighterServer) GetFile(context.Context, *FileRequest) (*FileReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetFile not implemented")
}
func (UnimplementedFreighterServer) GetTree(context.Context, *TreeRequest) (*TreeReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetTree not implemented")
}
func (UnimplementedFreighterServer) mustEmbedUnimplementedFreighterServer() {}

// UnsafeFreighterServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to FreighterServer will
// result in compilation errors.
type UnsafeFreighterServer interface {
	mustEmbedUnimplementedFreighterServer()
}

func RegisterFreighterServer(s grpc.ServiceRegistrar, srv FreighterServer) {
	s.RegisterService(&Freighter_ServiceDesc, srv)
}

func _Freighter_GetDir_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DirRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FreighterServer).GetDir(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Freighter_GetDir_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FreighterServer).GetDir(ctx, req.(*DirRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Freighter_GetFile_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(FileRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FreighterServer).GetFile(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Freighter_GetFile_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FreighterServer).GetFile(ctx, req.(*FileRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Freighter_GetTree_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(TreeRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FreighterServer).GetTree(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Freighter_GetTree_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FreighterServer).GetTree(ctx, req.(*TreeRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Freighter_ServiceDesc is the grpc.ServiceDesc for Freighter service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Freighter_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "freighter.Freighter",
	HandlerType: (*FreighterServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetDir",
			Handler:    _Freighter_GetDir_Handler,
		},
		{
			MethodName: "GetFile",
			Handler:    _Freighter_GetFile_Handler,
		},
		{
			MethodName: "GetTree",
			Handler:    _Freighter_GetTree_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/freighter.proto",
}
