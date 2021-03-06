// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package servicepb

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion7

// FabulaClient is the client API for Fabula service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type FabulaClient interface {
	WeaveSummary(ctx context.Context, in *WeaveSummaryRequest, opts ...grpc.CallOption) (*WeaveSummaryResponse, error)
}

type fabulaClient struct {
	cc grpc.ClientConnInterface
}

func NewFabulaClient(cc grpc.ClientConnInterface) FabulaClient {
	return &fabulaClient{cc}
}

var fabulaWeaveSummaryStreamDesc = &grpc.StreamDesc{
	StreamName: "WeaveSummary",
}

func (c *fabulaClient) WeaveSummary(ctx context.Context, in *WeaveSummaryRequest, opts ...grpc.CallOption) (*WeaveSummaryResponse, error) {
	out := new(WeaveSummaryResponse)
	err := c.cc.Invoke(ctx, "/merkleweave.protobuf.Fabula/WeaveSummary", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// FabulaService is the service API for Fabula service.
// Fields should be assigned to their respective handler implementations only before
// RegisterFabulaService is called.  Any unassigned fields will result in the
// handler for that method returning an Unimplemented error.
type FabulaService struct {
	WeaveSummary func(context.Context, *WeaveSummaryRequest) (*WeaveSummaryResponse, error)
}

func (s *FabulaService) weaveSummary(_ interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(WeaveSummaryRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return s.WeaveSummary(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     s,
		FullMethod: "/merkleweave.protobuf.Fabula/WeaveSummary",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return s.WeaveSummary(ctx, req.(*WeaveSummaryRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// RegisterFabulaService registers a service implementation with a gRPC server.
func RegisterFabulaService(s grpc.ServiceRegistrar, srv *FabulaService) {
	srvCopy := *srv
	if srvCopy.WeaveSummary == nil {
		srvCopy.WeaveSummary = func(context.Context, *WeaveSummaryRequest) (*WeaveSummaryResponse, error) {
			return nil, status.Errorf(codes.Unimplemented, "method WeaveSummary not implemented")
		}
	}
	sd := grpc.ServiceDesc{
		ServiceName: "merkleweave.protobuf.Fabula",
		Methods: []grpc.MethodDesc{
			{
				MethodName: "WeaveSummary",
				Handler:    srvCopy.weaveSummary,
			},
		},
		Streams:  []grpc.StreamDesc{},
		Metadata: "service.proto",
	}

	s.RegisterService(&sd, nil)
}
