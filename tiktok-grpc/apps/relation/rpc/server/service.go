package server

import (
	"context"
	"strconv"
	"tiktok-grpc/apps/relation/relation_pb"
	"tiktok-grpc/cmd"
	"tiktok-grpc/common/jwtx"
	"tiktok-grpc/common/utils"
	"tiktok-grpc/model"

	"gorm.io/gorm"
)

type RelationServiceServer struct {
	relation_pb.UnimplementedRelationServiceServer
}

func (p *RelationServiceServer) RelationAction(ctx context.Context, req *relation_pb.RelationActionRequest) (*relation_pb.RelationActionResponse, error) {
	token := req.Token
	toUserIDStr := req.ToUserId
	toUserID, err := strconv.Atoi(toUserIDStr)
	if err != nil {
		return &relation_pb.RelationActionResponse{
			StatusCode: 1, StatusMsg: "ID无效",
		}, nil

	}
	actionType := req.ActionType

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

	db := cmd.DB

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

			if err := cmd.RedisClient.Del(ctx, "user:"+useridstring).Err(); err != nil {
				return nil, err
			}
			if err := cmd.RedisClient.Del(ctx, "user:"+touseridstring).Err(); err != nil {
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
			if err := cmd.RedisClient.Del(ctx, "user:"+useridstring).Err(); err != nil {
				return nil, err
			}
			if err := cmd.RedisClient.Del(ctx, "user:"+touseridstring).Err(); err != nil {
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
func (p *RelationServiceServer) GetFollowList(ctx context.Context, req *relation_pb.GetFollowListRequest) (*relation_pb.GetFollowListResponse, error) {
	userID := req.UserId
	userIDInt, _ := strconv.Atoi(userID)

	token := req.Token

	tokenMsg, ok := jwtx.ParseToken(token)
	if int64(userIDInt) != tokenMsg.UserID && !ok {
		return &relation_pb.GetFollowListResponse{
			StatusCode: 1, StatusMsg: "无效token",
		}, nil
	}

	var follows []model.User
	db := cmd.DB

	err := db.Table("relation").
		Select("user.*").
		Joins("JOIN user ON user.id = relation.to_user_id").
		Where("relation.user_id = ? AND relation.is_follow=?", userID, true).
		Find(&follows).
		Error

	if err != nil {
		return &relation_pb.GetFollowListResponse{
			StatusCode: 1, StatusMsg: "查找错误",
		}, nil
	}
	followsMsg := []*relation_pb.UserInfo{}
	for i := range follows {

		follows[i].IsFollow = true
		followsMsg = append(followsMsg, utils.ConvertUserToRelationProto(follows[i]))
	}

	return &relation_pb.GetFollowListResponse{
		StatusCode: 0, StatusMsg: "关注列表",
		UserList: followsMsg,
	}, nil
}
func (p *RelationServiceServer) GetFansList(ctx context.Context, req *relation_pb.GetFansListRequest) (*relation_pb.GetFansListResponse, error) {
	userID := req.UserId
	userIDInt, _ := strconv.Atoi(userID)

	token := req.Token

	tokenMsg, ok := jwtx.ParseToken(token)
	if int64(userIDInt) != tokenMsg.UserID && !ok {

		return &relation_pb.GetFansListResponse{
			StatusCode: 1, StatusMsg: "无效token",
		}, nil
	}

	var fans []model.User
	db := cmd.DB

	err := db.Table("relation").
		Select("user.*").
		Joins("JOIN user ON user.id = relation.user_id").
		Where("relation.to_user_id = ? AND is_follow=?", userID, true).
		Find(&fans).
		Error

	if err != nil {
		return &relation_pb.GetFansListResponse{
			StatusCode: 1, StatusMsg: "查找错误",
		}, nil
	}

	var count int64

	fansMsg := []*relation_pb.UserInfo{}

	for i := range fans {

		db.Model(&model.Relation{}).Where("user_id = ? AND to_user_id = ? AND is_follow=?", userID, fans[i].ID, true).Count(&count)
		if count > 0 {
			fans[i].IsFollow = true
		}
		fansMsg = append(fansMsg, utils.ConvertUserToRelationProto(fans[i]))
	}

	return &relation_pb.GetFansListResponse{
		StatusCode: 0, StatusMsg: "粉丝列表",
		UserList: fansMsg,
	}, nil

}
func (p *RelationServiceServer) GetFriendList(ctx context.Context, req *relation_pb.GetFriendListRequest) (*relation_pb.GetFriendListResponse, error) {
	userIDStr := req.UserId
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {

		return &relation_pb.GetFriendListResponse{
			StatusCode: 1, StatusMsg: "无效ID",
		}, nil
	}

	token := req.Token
	tokenMsg, ok := jwtx.ParseToken(token)
	if tokenMsg.UserID != int64(userID) && !ok {

		return &relation_pb.GetFriendListResponse{
			StatusCode: 1, StatusMsg: "无效token",
		}, nil
	}

	var friends []model.User

	db := cmd.DB

	query := `
		SELECT user.*
		FROM relation 
		INNER JOIN user  ON relation .to_user_id = user.id
		WHERE relation.user_id = ? AND EXISTS (
			SELECT 1
			FROM relation
			WHERE user_id = user.id AND to_user_id = ?
		)
	`

	db.Raw(query, userID, userID).Scan(&friends)
	friendsMsg := []*relation_pb.UserInfo{}
	for i := range friends {
		friends[i].IsFollow = true
		friends[i].Avatar = cmd.StaticUrl + friends[i].Avatar
		friendsMsg = append(friendsMsg, utils.ConvertUserToRelationProto(friends[i]))
	}

	return &relation_pb.GetFriendListResponse{
		StatusCode: 0, StatusMsg: "好友列表",
		UserList: friendsMsg,
	}, nil

}
