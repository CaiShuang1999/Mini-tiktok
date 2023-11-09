package logic

import (
	"context"
	"strconv"

	"tiktok-micro/app/services/relation/internal/svc"
	"tiktok-micro/app/services/relation/relation_pb"
	"tiktok-micro/common/jwtx"
	"tiktok-micro/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetFriendListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetFriendListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetFriendListLogic {
	return &GetFriendListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetFriendListLogic) GetFriendList(in *relation_pb.GetFriendListRequest) (*relation_pb.GetFriendListResponse, error) {
	userIDStr := in.UserId
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {

		return &relation_pb.GetFriendListResponse{
			StatusCode: 1, StatusMsg: "无效ID",
		}, nil
	}

	token := in.Token
	tokenMsg, ok := jwtx.ParseToken(token)
	if tokenMsg.UserID != int64(userID) && !ok {

		return &relation_pb.GetFriendListResponse{
			StatusCode: 1, StatusMsg: "无效token",
		}, nil
	}

	var friends []model.User

	db := l.svcCtx.DB
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
		friends[i].Avatar = "http://" + l.svcCtx.Config.Nginx.Addr + friends[i].Avatar
		friendsMsg = append(friendsMsg, &relation_pb.UserInfo{
			Id:              friends[i].ID,
			Name:            friends[i].Name,
			Avatar:          friends[i].Avatar,
			BackgroundImage: friends[i].BackgroundImage,
			FavoriteCount:   friends[i].FavoriteCount,
			WorkCount:       friends[i].WorkCount,
			FollowCount:     friends[i].FollowCount,
			FollowerCount:   friends[i].FollowerCount,
			Signature:       friends[i].Signature,
			TotalFavorited:  friends[i].TotalFavorited,
			IsFollow:        true,
		})
	}

	return &relation_pb.GetFriendListResponse{
		StatusCode: 0, StatusMsg: "好友列表",
		UserList: friendsMsg,
	}, nil
}
