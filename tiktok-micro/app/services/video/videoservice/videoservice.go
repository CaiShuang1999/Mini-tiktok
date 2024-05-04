// Code generated by goctl. DO NOT EDIT.
// Source: video.proto

package videoservice

import (
	"context"

	"tiktok-micro/app/services/video/video_pb"

	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
)

type (
	GetVideoInfoRequest  = video_pb.GetVideoInfoRequest
	GetVideoInfoResponse = video_pb.GetVideoInfoResponse
	PublishlistRequest   = video_pb.PublishlistRequest
	PublishlistResponse  = video_pb.PublishlistResponse
	UploadVideoRequest   = video_pb.UploadVideoRequest
	UploadVideoResponse  = video_pb.UploadVideoResponse
	UserInfo             = video_pb.UserInfo
	VideoInfo            = video_pb.VideoInfo

	VideoService interface {
		UploadVideo(ctx context.Context, in *UploadVideoRequest, opts ...grpc.CallOption) (*UploadVideoResponse, error)
		PublishList(ctx context.Context, in *PublishlistRequest, opts ...grpc.CallOption) (*PublishlistResponse, error)
		GetVideoInfo(ctx context.Context, in *GetVideoInfoRequest, opts ...grpc.CallOption) (*GetVideoInfoResponse, error)
	}

	defaultVideoService struct {
		cli zrpc.Client
	}
)

func NewVideoService(cli zrpc.Client) VideoService {
	return &defaultVideoService{
		cli: cli,
	}
}

func (m *defaultVideoService) UploadVideo(ctx context.Context, in *UploadVideoRequest, opts ...grpc.CallOption) (*UploadVideoResponse, error) {
	client := video_pb.NewVideoServiceClient(m.cli.Conn())
	return client.UploadVideo(ctx, in, opts...)
}

func (m *defaultVideoService) PublishList(ctx context.Context, in *PublishlistRequest, opts ...grpc.CallOption) (*PublishlistResponse, error) {
	client := video_pb.NewVideoServiceClient(m.cli.Conn())
	return client.PublishList(ctx, in, opts...)
}

func (m *defaultVideoService) GetVideoInfo(ctx context.Context, in *GetVideoInfoRequest, opts ...grpc.CallOption) (*GetVideoInfoResponse, error) {
	client := video_pb.NewVideoServiceClient(m.cli.Conn())
	return client.GetVideoInfo(ctx, in, opts...)
}