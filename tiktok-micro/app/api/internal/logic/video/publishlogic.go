package video

import (
	"context"
	"io"
	"net/http"
	"tiktok-micro/app/api/internal/svc"
	"tiktok-micro/app/api/internal/types"
	"tiktok-micro/app/services/video/video_pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type PublishLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPublishLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PublishLogic {
	return &PublishLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PublishLogic) Publish(req *types.PublishRequest, r *http.Request) (resp *types.PublishResponse, err error) {
	file, _, err := r.FormFile("data")
	if err != nil {
		return nil, err
	}
	defer file.Close()
	videoContent, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	res, err := l.svcCtx.VideoRpc.UploadVideo(
		l.ctx, &video_pb.UploadVideoRequest{
			Data:  []byte(videoContent),
			Token: req.Token,
			Title: req.Title,
		})
	if err != nil {
		return nil, err
	}

	return &types.PublishResponse{
		StatusCode: res.StatusCode,
		StatusMsg:  res.StatusMsg,
	}, nil
}
