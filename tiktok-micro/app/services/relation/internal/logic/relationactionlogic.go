package logic

import (
	"context"
	"strconv"

	"tiktok-micro/app/services/relation/internal/svc"
	"tiktok-micro/app/services/relation/relation_pb"
	"tiktok-micro/common/jwtx"
	"tiktok-micro/model"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type RelationActionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRelationActionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RelationActionLogic {
	return &RelationActionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *RelationActionLogic) RelationAction(in *relation_pb.RelationActionRequest) (*relation_pb.RelationActionResponse, error) {
	ctx := context.Background()
	token := in.Token
	toUserIDStr := in.ToUserId
	toUserID, err := strconv.Atoi(toUserIDStr)
	if err != nil {
		return &relation_pb.RelationActionResponse{
			StatusCode: 1, StatusMsg: "ID无效",
		}, nil

	}
	actionType := in.ActionType

	tokenMsg, ok := jwtx.ParseToken(token)

	if !ok {
		return &relation_pb.RelationActionResponse{
			StatusCode: 1, StatusMsg: "token无效",
		}, nil
	}

	useridstring := strconv.Itoa(int(tokenMsg.UserID))

	var relation model.Relation
	relation.UserID = tokenMsg.UserID
	relation.ToUserID = int64(toUserID)
	if relation.UserID == relation.ToUserID {
		return &relation_pb.RelationActionResponse{
			StatusCode: 1, StatusMsg: "不能对自己操作!",
		}, nil

	}

	db := l.svcCtx.DB

	if db.Where("user_id=? AND to_user_id=?", relation.UserID, relation.ToUserID).Find(&model.Relation{}).RowsAffected == 0 {
		db.Create(&relation)
	}

	touseridstring := strconv.Itoa(int(relation.ToUserID))
	switch actionType {
	case "1":
		{

			db.Model(&model.Relation{}).Where("user_id=? AND to_user_id=?", relation.UserID, relation.ToUserID).
				Update("is_follow", true)
			db.Model(&model.User{}).Where("id = ?", relation.UserID).Update("follow_count", gorm.Expr("follow_count + ?", 1))
			db.Model(&model.User{}).Where("id = ?", relation.ToUserID).Update("follower_count", gorm.Expr("follower_count + ?", 1))

			if err := l.svcCtx.Rdb_0.Del(ctx, "user:"+useridstring).Err(); err != nil {
				return nil, err
			}
			if err := l.svcCtx.Rdb_0.Del(ctx, "user:"+touseridstring).Err(); err != nil {
				return nil, err
			}

			return &relation_pb.RelationActionResponse{
				StatusCode: 0, StatusMsg: "关注成功!",
			}, nil

		}

	case "2":
		{
			db.Model(&model.Relation{}).Where("user_id=? AND to_user_id=?", relation.UserID, relation.ToUserID).
				Update("is_follow", false)
			db.Model(&model.User{}).Where("id = ?", relation.UserID).Update("follow_count", gorm.Expr("follow_count - ?", 1))
			db.Model(&model.User{}).Where("id = ?", relation.ToUserID).Update("follower_count", gorm.Expr("follower_count - ?", 1))
			if err := l.svcCtx.Rdb_0.Del(ctx, "user:"+useridstring).Err(); err != nil {
				return nil, err
			}
			if err := l.svcCtx.Rdb_0.Del(ctx, "user:"+touseridstring).Err(); err != nil {
				return nil, err
			}
			return &relation_pb.RelationActionResponse{
				StatusCode: 0, StatusMsg: "取消关注成功!",
			}, nil
		}
	default:
		{
			return &relation_pb.RelationActionResponse{
				StatusCode: 1, StatusMsg: "无效操作!",
			}, nil

		}

	}
}
