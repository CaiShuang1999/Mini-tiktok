package logic

import (
	"context"
	"log"
	"os"
	"os/exec"
	"strconv"
	"time"

	"tiktok-micro/app/services/video/internal/svc"
	"tiktok-micro/app/services/video/video_pb"
	"tiktok-micro/common/jwtx"
	"tiktok-micro/model"

	"github.com/redis/go-redis/v9"
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

	tokenmsg, ok := jwtx.ParseToken(in.Token)
	if !ok {

		return &video_pb.UploadVideoResponse{
			StatusCode: 0, StatusMsg: "无效token",
		}, nil
	}

	useridstring := strconv.Itoa(int(tokenmsg.UserID))
	timestamp := time.Now().UnixNano() / int64(time.Millisecond) // 获取当前时间的时间戳
	timestampString := strconv.Itoa(int(timestamp))
	fileNameType := "mp4"
	savePath := "../../../web/static/video/" + useridstring + "_" + timestampString + "." + fileNameType // 保存路径（用户ID+上传时间）

	err := os.WriteFile(savePath, in.Data, 0644)

	if err != nil {
		log.Printf("Failed to save video file: %v", err)
		return nil, err
	}
	// 封面
	inputFile := savePath
	outputFile := "../../../web/static/image/" + useridstring + "_" + timestampString + ".jpg"
	timeOffset := "00:00:01"

	cmd := exec.Command("ffmpeg", "-i", inputFile, "-ss", timeOffset, "-vframes", "1", outputFile)

	err = cmd.Run()

	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	// 上传成功后，video表数据库
	db := l.svcCtx.DB
	var newVideo model.Video
	newVideo.UserID = tokenmsg.UserID
	newVideo.PlayURL = "http://" + l.svcCtx.Config.Nginx.Addr + "/video/" + useridstring + "_" + timestampString + "." + fileNameType
	newVideo.Title = in.Title
	newVideo.CoverURL = "http://" + l.svcCtx.Config.Nginx.Addr + "/image/" + useridstring + "_" + timestampString + ".jpg"
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
	err = l.svcCtx.Rdb_0.Del(context.Background(), "user:"+useridstring).Err()
	if err != nil {

		return nil, err
	}

	_, err = VideoCache(l.svcCtx.DB, l.svcCtx.Rdb_0, l.svcCtx.Rdb_1, newVideo.ID)

	if err != nil {
		return nil, err
	}

	return &video_pb.UploadVideoResponse{StatusCode: 0, StatusMsg: "上传成功!"}, nil
}

var ctx = context.Background()

func CreateVideoInfoCache(db *gorm.DB, rdb *redis.Client, video_id string) (model.Video, error) {
	var video model.Video

	hashKey := "video:" + video_id

	result := db.First(&video, video_id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return model.Video{}, result.Error
		} else {

			return model.Video{}, result.Error
		}
	}

	if err := rdb.HSet(ctx, hashKey, video).Err(); err != nil {
		return model.Video{}, err
	}

	return video, nil
}

func VideoInfoCache(db *gorm.DB, rdb *redis.Client, video_id string) (model.Video, error) {
	var video model.Video

	hashKey := "video:" + video_id //视频key
	video_id_int, _ := strconv.Atoi(video_id)
	exists, err := rdb.Exists(ctx, hashKey).Result() //用户信息是否存在缓存
	if err != nil {
		return model.Video{}, err
	}

	if exists == 1 {
		//用户缓存加载
		if err := rdb.HGetAll(ctx, hashKey).Scan(&video); err != nil {
			return model.Video{}, err
		}

	} else {
		video, err = CreateVideoInfoCache(db, rdb, video_id)
	}

	video.ID = int64(video_id_int)

	return video, err

}

func CreateVideoFavoritedCache(db *gorm.DB, rdb2 *redis.Client, video_id string) (int64, error) {
	var video model.Video

	var count int64 = 0

	result := db.First(&video, video_id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return count, result.Error
		} else {

			return count, result.Error
		}
	}

	videoFavoriteCountKey := "favorite:video:" + video_id
	count = video.FavoriteCount

	rdb2.Set(ctx, videoFavoriteCountKey, count, 0)

	return count, nil
}

func VideoFavoritedCache(db *gorm.DB, rdb2 *redis.Client, video_id string) (int64, error) {

	var count int64 = 0
	videoFavoriteCountKey := "favorite:video:" + video_id            //视频点赞数key
	exists2, err := rdb2.Exists(ctx, videoFavoriteCountKey).Result() //视频点赞数key是否存在

	if err != nil {
		return count, err
	}

	if exists2 == 1 {
		//点赞数缓存加载

		videoFavoriteCountStr, _ := rdb2.Get(ctx, videoFavoriteCountKey).Result()
		videoFavoriteCount, err := strconv.Atoi(videoFavoriteCountStr)
		if err != nil {
			return count, err
		}
		count = int64(videoFavoriteCount)
	} else {
		count, err = CreateVideoFavoritedCache(db, rdb2, video_id)

	}
	return count, err

}

func VideoCache(db *gorm.DB, rdb *redis.Client, rdb2 *redis.Client, videoID int64) (model.Video, error) {

	video_id := strconv.Itoa(int(videoID))
	video, err := VideoInfoCache(db, rdb, video_id)
	if err != nil {
		return video, err
	}
	videoFavoriteCount, err := VideoFavoritedCache(db, rdb2, video_id)
	if err != nil {
		return video, err
	}

	video.FavoriteCount = videoFavoriteCount

	return video, err
}
