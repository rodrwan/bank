// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package users

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion7

// UsersReadServiceClient is the client API for UsersReadService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type UsersReadServiceClient interface {
	// Article returns a single article by ID
	GetUser(ctx context.Context, in *GetUserRequest, opts ...grpc.CallOption) (*GetUserResponse, error)
}

type usersReadServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewUsersReadServiceClient(cc grpc.ClientConnInterface) UsersReadServiceClient {
	return &usersReadServiceClient{cc}
}

func (c *usersReadServiceClient) GetUser(ctx context.Context, in *GetUserRequest, opts ...grpc.CallOption) (*GetUserResponse, error) {
	out := new(GetUserResponse)
	err := c.cc.Invoke(ctx, "/users.UsersReadService/GetUser", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// UsersReadServiceServer is the server API for UsersReadService service.
// All implementations must embed UnimplementedUsersReadServiceServer
// for forward compatibility
type UsersReadServiceServer interface {
	// Article returns a single article by ID
	GetUser(context.Context, *GetUserRequest) (*GetUserResponse, error)
	mustEmbedUnimplementedUsersReadServiceServer()
}

// UnimplementedUsersReadServiceServer must be embedded to have forward compatible implementations.
type UnimplementedUsersReadServiceServer struct {
}

func (UnimplementedUsersReadServiceServer) GetUser(context.Context, *GetUserRequest) (*GetUserResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetUser not implemented")
}
func (UnimplementedUsersReadServiceServer) mustEmbedUnimplementedUsersReadServiceServer() {}

// UnsafeUsersReadServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to UsersReadServiceServer will
// result in compilation errors.
type UnsafeUsersReadServiceServer interface {
	mustEmbedUnimplementedUsersReadServiceServer()
}

func RegisterUsersReadServiceServer(s grpc.ServiceRegistrar, srv UsersReadServiceServer) {
	s.RegisterService(&_UsersReadService_serviceDesc, srv)
}

func _UsersReadService_GetUser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetUserRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UsersReadServiceServer).GetUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/users.UsersReadService/GetUser",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UsersReadServiceServer).GetUser(ctx, req.(*GetUserRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _UsersReadService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "users.UsersReadService",
	HandlerType: (*UsersReadServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetUser",
			Handler:    _UsersReadService_GetUser_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "users.proto",
}

// UsersWriteServiceClient is the client API for UsersWriteService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type UsersWriteServiceClient interface {
	// CreateUser creates a User.
	CreateUser(ctx context.Context, in *CreateUserRequest, opts ...grpc.CallOption) (*CreateUserResponse, error)
}

type usersWriteServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewUsersWriteServiceClient(cc grpc.ClientConnInterface) UsersWriteServiceClient {
	return &usersWriteServiceClient{cc}
}

func (c *usersWriteServiceClient) CreateUser(ctx context.Context, in *CreateUserRequest, opts ...grpc.CallOption) (*CreateUserResponse, error) {
	out := new(CreateUserResponse)
	err := c.cc.Invoke(ctx, "/users.UsersWriteService/CreateUser", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// UsersWriteServiceServer is the server API for UsersWriteService service.
// All implementations must embed UnimplementedUsersWriteServiceServer
// for forward compatibility
type UsersWriteServiceServer interface {
	// CreateUser creates a User.
	CreateUser(context.Context, *CreateUserRequest) (*CreateUserResponse, error)
	mustEmbedUnimplementedUsersWriteServiceServer()
}

// UnimplementedUsersWriteServiceServer must be embedded to have forward compatible implementations.
type UnimplementedUsersWriteServiceServer struct {
}

func (UnimplementedUsersWriteServiceServer) CreateUser(context.Context, *CreateUserRequest) (*CreateUserResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateUser not implemented")
}
func (UnimplementedUsersWriteServiceServer) mustEmbedUnimplementedUsersWriteServiceServer() {}

// UnsafeUsersWriteServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to UsersWriteServiceServer will
// result in compilation errors.
type UnsafeUsersWriteServiceServer interface {
	mustEmbedUnimplementedUsersWriteServiceServer()
}

func RegisterUsersWriteServiceServer(s grpc.ServiceRegistrar, srv UsersWriteServiceServer) {
	s.RegisterService(&_UsersWriteService_serviceDesc, srv)
}

func _UsersWriteService_CreateUser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateUserRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UsersWriteServiceServer).CreateUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/users.UsersWriteService/CreateUser",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UsersWriteServiceServer).CreateUser(ctx, req.(*CreateUserRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _UsersWriteService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "users.UsersWriteService",
	HandlerType: (*UsersWriteServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateUser",
			Handler:    _UsersWriteService_CreateUser_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "users.proto",
}