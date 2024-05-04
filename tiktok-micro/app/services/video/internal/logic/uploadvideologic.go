package logic

import (
	"context"
	"strconv"

	"tiktok-micro/app/services/video/internal/svc"
	"tiktok-micro/app/services/video/video_pb"
	"tiktok-micro/common/jwtx"
	"tiktok-micro/model"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type UploadVideoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUploadVideoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UploadVideoLogic {
	return &UploadVideoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UploadVideoLogic) UploadVideo(in *video_pb.UploadVideoRequest) (*video_pb.UploadVideoResponse, error) {
	// todo: add your logic here and delete this line

	timestamp := in.Timestamp
	timestampString := strconv.Itoa(int(timestamp))

	tokenmsg, ok := jwtx.ParseToken(in.Token)
	if !ok {

		return &video_pb.UploadVideoResponse{
			StatusCode: 0, StatusMsg: "无效token",
		}, nil
	}
	useridstring := strconv.Itoa(int(tokenmsg.UserID))
	// 上传成功后，video表数据库
	db := l.svcCtx.DB

	var newVideo model.Video
	newVideo.UserID = tokenmsg.UserID
	newVideo.PlayURL = "http://" + l.svcCtx.Config.Nginx.Addr + "/video/" + useridstring + "_" + timestampString + ".mp4"
	newVideo.Title = in.Title
	newVideo.CoverURL = "http://" + l.svcCtx.Config.Nginx.Addr + "/image/" + useridstring + "_" + timestampString + ".jpg"
	newVideo.CreateTime = timestamp

	//开启事务
	tx := db.Begin()
	if err := tx.Create(&newVideo).Error; err != nil { //创建video记录
		tx.Rollback()

		return &video_pb.UploadVideoResponse{StatusCode: 1, StatusMsg: "创建数据库记录失败！"}, nil
	}

	//sql记录修改(作品+1)
	if err := tx.Model(&model.User{}).Where("id = ?", tokenmsg.UserID).Update("work_count", gorm.Expr("work_count + ?", 1)).Error; err != nil {
		tx.Rollback()
		return &video_pb.UploadVideoResponse{StatusCode: 1, StatusMsg: "创建数据库记录失败！"}, nil
	}

	tx.Commit() //提交事务

	//缓存删除（一致性）
	err := l.svcCtx.Rdb_0.Del(context.Background(), "user:"+useridstring).Err()
	if err != nil {

		return nil, err
	}

	_, err = VideoCache(l.svcCtx.DB, l.svcCtx.Rdb_0, l.svcCtx.Rdb_1, newVideo.ID)

	if err != nil {
		return nil, err
	}

	return &video_pb.UploadVideoResponse{StatusCode: 0, StatusMsg: "上传成功!"}, nil
}
