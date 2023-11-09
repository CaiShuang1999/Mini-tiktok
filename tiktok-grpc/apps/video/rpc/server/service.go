package server

import (
	"context"
	"log"
	"os"
	"os/exec"

	"strconv"
	"tiktok-grpc/apps/video/video_pb"
	cmdx "tiktok-grpc/cmd"
	"tiktok-grpc/common/jwtx"
	"tiktok-grpc/common/redisx"
	"tiktok-grpc/common/utils"
	"tiktok-grpc/model"
	"time"

	"gorm.io/gorm"
)

type VideoServiceServer struct {
	video_pb.UnimplementedVideoServiceServer
}

func (p *VideoServiceServer) UploadVideo(ctx context.Context, req *video_pb.UploadVideoRequest) (*video_pb.UploadVideoResponse, error) {
	tokenmsg, ok := jwtx.ParseToken(req.Token)
	if !ok {

		return &video_pb.UploadVideoResponse{
			StatusCode: 0, StatusMsg: "无效token",
		}, nil
	}

	useridstring := strconv.Itoa(int(tokenmsg.UserID))
	timestamp := time.Now().UnixNano() / int64(time.Millisecond) // 获取当前时间的时间戳
	timestampString := strconv.Itoa(int(timestamp))
	fileNameType := "mp4"
	savePath := "web/static/video/" + useridstring + "_" + timestampString + "." + fileNameType // 保存路径（用户ID+上传时间）
	err := os.WriteFile(savePath, req.Data, 0644)
	if err != nil {
		log.Printf("Failed to save video file: %v", err)
		return nil, err
	}
	//封面
	inputFile := savePath
	outputFile := "web/static/image/" + useridstring + "_" + timestampString + ".jpg"
	timeOffset := "00:00:01"

	cmd := exec.Command("ffmpeg", "-i", inputFile, "-ss", timeOffset, "-vframes", "1", outputFile)

	err = cmd.Run()

	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	//上传成功后，video表数据库
	db := cmdx.DB
	var newVideo model.Video
	newVideo.UserID = tokenmsg.UserID
	newVideo.PlayURL = cmdx.StaticUrl + "/video/" + useridstring + "_" + timestampString + "." + fileNameType
	newVideo.Title = req.Title
	newVideo.CoverURL = cmdx.StaticUrl + "/image/" + useridstring + "_" + timestampString + ".jpg"
	newVideo.CreateTime = timestamp
	result := db.Create(&newVideo) //数据库创建video记录
	if result.Error != nil {
		os.Remove(outputFile) //删除已上传的视频，保持一致性
		os.Remove(savePath)
		return &video_pb.UploadVideoResponse{StatusCode: 1, StatusMsg: "创建数据库记录失败！"}, nil
	}

	//sql记录修改(作品+1)
	db.Model(&model.User{}).Where("id = ?", tokenmsg.UserID).Update("work_count", gorm.Expr("work_count + ?", 1))

	//缓存删除（一致性）
	err = cmdx.RedisClient.Del(ctx, "user:"+useridstring).Err()
	if err != nil {

		return nil, err
	}

	_, err = redisx.VideoCache(newVideo.ID)
	if err != nil {
		return nil, err
	}

	return &video_pb.UploadVideoResponse{StatusCode: 0, StatusMsg: "上传成功!"}, nil
}
func (p *VideoServiceServer) PublishList(ctx context.Context, req *video_pb.PublishlistRequest) (*video_pb.PublishlistResponse, error) {

	userID := req.UserId
	var videos []model.Video
	db := cmdx.DB

	db.Preload("Author").Where("user_id = ?", userID).Find(&videos) //使用db.Preload("Author")预加载user表中"Author"信息，否则为0值

	videos_msg := []*video_pb.VideoInfo{}
	for i := range videos {
		videos_msg = append(videos_msg, utils.ConvertVideoToProto(videos[i]))
	}

	return &video_pb.PublishlistResponse{
		StatusCode: 0,
		StatusMsg:  "发布列表",
		VideoList:  videos_msg,
	}, nil

}
