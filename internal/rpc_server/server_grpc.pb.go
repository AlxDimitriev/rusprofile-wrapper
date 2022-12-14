// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             (unknown)
// source: internal/rpc_server/server.proto

package rpc_server

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

// CompanyInfoServiceClient is the client API for CompanyInfoService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type CompanyInfoServiceClient interface {
	FetchCompanyInfo(ctx context.Context, in *CompanyRequest, opts ...grpc.CallOption) (*CompanyResponse, error)
}

type companyInfoServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewCompanyInfoServiceClient(cc grpc.ClientConnInterface) CompanyInfoServiceClient {
	return &companyInfoServiceClient{cc}
}

func (c *companyInfoServiceClient) FetchCompanyInfo(ctx context.Context, in *CompanyRequest, opts ...grpc.CallOption) (*CompanyResponse, error) {
	out := new(CompanyResponse)
	err := c.cc.Invoke(ctx, "/rpc_server.CompanyInfoService/FetchCompanyInfo", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// CompanyInfoServiceServer is the server API for CompanyInfoService service.
// All implementations must embed UnimplementedCompanyInfoServiceServer
// for forward compatibility
type CompanyInfoServiceServer interface {
	FetchCompanyInfo(context.Context, *CompanyRequest) (*CompanyResponse, error)
	mustEmbedUnimplementedCompanyInfoServiceServer()
}

// UnimplementedCompanyInfoServiceServer must be embedded to have forward compatible implementations.
type UnimplementedCompanyInfoServiceServer struct {
}

func (UnimplementedCompanyInfoServiceServer) FetchCompanyInfo(context.Context, *CompanyRequest) (*CompanyResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method FetchCompanyInfo not implemented")
}
func (UnimplementedCompanyInfoServiceServer) mustEmbedUnimplementedCompanyInfoServiceServer() {}

// UnsafeCompanyInfoServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to CompanyInfoServiceServer will
// result in compilation errors.
type UnsafeCompanyInfoServiceServer interface {
	mustEmbedUnimplementedCompanyInfoServiceServer()
}

func RegisterCompanyInfoServiceServer(s grpc.ServiceRegistrar, srv CompanyInfoServiceServer) {
	s.RegisterService(&CompanyInfoService_ServiceDesc, srv)
}

func _CompanyInfoService_FetchCompanyInfo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CompanyRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CompanyInfoServiceServer).FetchCompanyInfo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/rpc_server.CompanyInfoService/FetchCompanyInfo",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CompanyInfoServiceServer).FetchCompanyInfo(ctx, req.(*CompanyRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// CompanyInfoService_ServiceDesc is the grpc.ServiceDesc for CompanyInfoService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var CompanyInfoService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "rpc_server.CompanyInfoService",
	HandlerType: (*CompanyInfoServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "FetchCompanyInfo",
			Handler:    _CompanyInfoService_FetchCompanyInfo_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "internal/rpc_server/server.proto",
}
