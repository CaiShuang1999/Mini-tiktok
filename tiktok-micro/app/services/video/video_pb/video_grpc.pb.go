// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.19.4
// source: video.proto

package video_pb

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

// VideoServiceClient is the client API for VideoService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type VideoServiceClient interface {
	UploadVideo(ctx context.Context, in *UploadVideoRequest, opts ...grpc.CallOption) (*UploadVideoResponse, error)
	PublishList(ctx context.Context, in *PublishlistRequest, opts ...grpc.CallOption) (*PublishlistResponse, error)
	GetVideoInfo(ctx context.Context, in *GetVideoInfoRequest, opts ...grpc.CallOption) (*GetVideoInfoResponse, error)
}

type videoServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewVideoServiceClient(cc grpc.ClientConnInterface) VideoServiceClient {
	return &videoServiceClient{cc}
}

func (c *videoServiceClient) UploadVideo(ctx context.Context, in *UploadVideoRequest, opts ...grpc.CallOption) (*UploadVideoResponse, error) {
	out := new(UploadVideoResponse)
	err := c.cc.Invoke(ctx, "/video_pb.VideoService/UploadVideo", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *videoServiceClient) PublishList(ctx context.Context, in *PublishlistRequest, opts ...grpc.CallOption) (*PublishlistResponse, error) {
	out := new(PublishlistResponse)
	err := c.cc.Invoke(ctx, "/video_pb.VideoService/PublishList", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *videoServiceClient) GetVideoInfo(ctx context.Context, in *GetVideoInfoRequest, opts ...grpc.CallOption) (*GetVideoInfoResponse, error) {
	out := new(GetVideoInfoResponse)
	err := c.cc.Invoke(ctx, "/video_pb.VideoService/GetVideoInfo", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// VideoServiceServer is the server API for VideoService service.
// All implementations must embed UnimplementedVideoServiceServer
// for forward compatibility
type VideoServiceServer interface {
	UploadVideo(context.Context, *UploadVideoRequest) (*UploadVideoResponse, error)
	PublishList(context.Context, *PublishlistRequest) (*PublishlistResponse, error)
	GetVideoInfo(context.Context, *GetVideoInfoRequest) (*GetVideoInfoResponse, error)
	mustEmbedUnimplementedVideoServiceServer()
}

// UnimplementedVideoServiceServer must be embedded to have forward compatible implementations.
type UnimplementedVideoServiceServer struct {
}

func (UnimplementedVideoServiceServer) UploadVideo(context.Context, *UploadVideoRequest) (*UploadVideoResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UploadVideo not implemented")
}
func (UnimplementedVideoServiceServer) PublishList(context.Context, *PublishlistRequest) (*PublishlistResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PublishList not implemented")
}
func (UnimplementedVideoServiceServer) GetVideoInfo(context.Context, *GetVideoInfoRequest) (*GetVideoInfoResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetVideoInfo not implemented")
}
func (UnimplementedVideoServiceServer) mustEmbedUnimplementedVideoServiceServer() {}

// UnsafeVideoServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to VideoServiceServer will
// result in compilation errors.
type UnsafeVideoServiceServer interface {
	mustEmbedUnimplementedVideoServiceServer()
}

func RegisterVideoServiceServer(s grpc.ServiceRegistrar, srv VideoServiceServer) {
	s.RegisterService(&VideoService_ServiceDesc, srv)
}

func _VideoService_UploadVideo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UploadVideoRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(VideoServiceServer).UploadVideo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/video_pb.VideoService/UploadVideo",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(VideoServiceServer).UploadVideo(ctx, req.(*UploadVideoRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _VideoService_PublishList_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PublishlistRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(VideoServiceServer).PublishList(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/video_pb.VideoService/PublishList",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(VideoServiceServer).PublishList(ctx, req.(*PublishlistRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _VideoService_GetVideoInfo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetVideoInfoRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(VideoServiceServer).GetVideoInfo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/video_pb.VideoService/GetVideoInfo",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(VideoServiceServer).GetVideoInfo(ctx, req.(*GetVideoInfoRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// VideoService_ServiceDesc is the grpc.ServiceDesc for VideoService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var VideoService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "video_pb.VideoService",
	HandlerType: (*VideoServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "UploadVideo",
			Handler:    _VideoService_UploadVideo_Handler,
		},
		{
			MethodName: "PublishList",
			Handler:    _VideoService_PublishList_Handler,
		},
		{
			MethodName: "GetVideoInfo",
			Handler:    _VideoService_GetVideoInfo_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "video.proto",
}
