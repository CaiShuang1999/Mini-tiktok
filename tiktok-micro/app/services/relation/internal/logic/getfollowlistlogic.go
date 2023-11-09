package logic

import (
	"context"
	"strconv"

	"tiktok-micro/app/services/relation/internal/svc"
	"tiktok-micro/app/services/relation/relation_pb"
	"tiktok-micro/app/services/user/user_pb"
	"tiktok-micro/common/jwtx"
	"tiktok-micro/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetFollowListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetFollowListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetFollowListLogic {
	return &GetFollowListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetFollowListLogic) GetFollowList(in *relation_pb.GetFollowListRequest) (*relation_pb.GetFollowListResponse, error) {
	userID := in.UserId
	userIDInt, _ := strconv.Atoi(userID)

	token := in.Token

	tokenMsg, ok := jwtx.ParseToken(token)
	if int64(userIDInt) != tokenMsg.UserID && !ok {
		return &relation_pb.GetFollowListResponse{
			StatusCode: 1, StatusMsg: "无效token",
		}, nil
	}

	var followsRelation []model.Relation
	db := l.svcCtx.DB

	err := db.Where("user_id = ? AND is_follow=?", userID, true).Find(&followsRelation).Error

	if err != nil {
		return &relation_pb.GetFollowListResponse{
			StatusCode: 1, StatusMsg: "查找错误",
		}, nil
	}
	followsMsg := []*relation_pb.UserInfo{}
	for i := range followsRelation {
		follows_idstr := strconv.Itoa(int(followsRelation[i].ToUserID))
		followsinfo, _ := l.svcCtx.UserRpc.GetUserInfo(context.Background(), &user_pb.UserInfoRequest{UserId: follows_idstr, Token: token})

		followsMsg = append(followsMsg, &relation_pb.UserInfo{
			Id:              followsinfo.User.Id,
			Name:            followsinfo.User.Name,
			Avatar:          followsinfo.User.Avatar,
			BackgroundImage: followsinfo.User.BackgroundImage,
			FavoriteCount:   followsinfo.User.FavoriteCount,
			FollowCount:     followsinfo.User.FollowCount,
			FollowerCount:   followsinfo.User.FollowerCount,
			WorkCount:       followsinfo.User.WorkCount,
			Signature:       followsinfo.User.Signature,
			TotalFavorited:  followsinfo.User.TotalFavorited,
			IsFollow:        true,
		})
	}

	return &relation_pb.GetFollowListResponse{
		StatusCode: 0, StatusMsg: "关注列表",
		UserList: followsMsg,
	}, nil
}
