package video

import (
	"context"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"tiktok-micro/app/api/internal/svc"
	"tiktok-micro/app/api/internal/types"
	"tiktok-micro/app/services/video/video_pb"
	"tiktok-micro/common/jwtx"
	"time"

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

	//保存视频
	tokenmsg, ok := jwtx.ParseToken(req.Token)
	if !ok {

		return &types.PublishResponse{
			StatusCode: 0, StatusMsg: "无效token",
		}, nil
	}
	useridstring := strconv.Itoa(int(tokenmsg.UserID))
	timestamp := time.Now().UnixNano() / int64(time.Millisecond) // 获取当前时间的时间戳
	timestampString := strconv.Itoa(int(timestamp))
	fileNameType := "mp4"

	savePath := "../../web/static/video/" + useridstring + "_" + timestampString + "." + fileNameType // 保存路径（用户ID+上传时间）
	err = os.WriteFile(savePath, videoContent, 0644)
	if err != nil {
		log.Printf("Failed to save video file: %v", err)
		return nil, err
	}

	// 封面
	inputFile := savePath
	outputFile := "../../web/static/image/" + useridstring + "_" + timestampString + ".jpg"
	timeOffset := "00:00:01"
	cmd := exec.Command("ffmpeg", "-i", inputFile, "-ss", timeOffset, "-vframes", "1", outputFile)
	err = cmd.Run()
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	//只执行数据库

	res, err := l.svcCtx.VideoRpc.UploadVideo(
		l.ctx, &video_pb.UploadVideoRequest{
			Timestamp: timestamp,
			Token:     req.Token,
			Title:     req.Title,
		})
	if res.StatusCode == 1 {
		os.Remove(outputFile)
		os.Remove(savePath)
	}
	if err != nil {
		os.Remove(outputFile)
		os.Remove(savePath)
		return nil, err
	}

	return &types.PublishResponse{
		StatusCode: res.StatusCode,
		StatusMsg:  res.StatusMsg,
	}, nil
}
