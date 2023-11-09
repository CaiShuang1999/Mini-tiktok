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

type GetFansListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetFansListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetFansListLogic {
	return &GetFansListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetFansListLogic) GetFansList(in *relation_pb.GetFansListRequest) (*relation_pb.GetFansListResponse, error) {
	// todo: add your logic here and delete this line

	userID := in.UserId
	userIDInt, _ := strconv.Atoi(userID)

	token := in.Token

	tokenMsg, ok := jwtx.ParseToken(token)
	if int64(userIDInt) != tokenMsg.UserID && !ok {

		return &relation_pb.GetFansListResponse{
			StatusCode: 1, StatusMsg: "无效token",
		}, nil
	}

	var fans []model.Relation
	db := l.svcCtx.DB

	err := db.Where("to_user_id = ? AND is_follow=?", userID, true).Find(&fans).Error
	if err != nil {
		return &relation_pb.GetFansListResponse{
			StatusCode: 1, StatusMsg: "查找错误",
		}, nil
	}

	var count int64

	fansMsg := []*relation_pb.UserInfo{}

	for i := range fans {
		fans_idstr := strconv.Itoa(int(fans[i].UserID))
		fansinfo, _ := l.svcCtx.UserRpc.GetUserInfo(context.Background(), &user_pb.UserInfoRequest{UserId: fans_idstr, Token: token})
		isFollow := false
		db.Model(&model.Relation{}).Where("user_id = ? AND to_user_id = ? AND is_follow=?", userID, fans[i].UserID, true).Count(&count)
		if count > 0 {
			isFollow = true
		}
		fansMsg = append(fansMsg, &relation_pb.UserInfo{
			Id:              fansinfo.User.Id,
			Name:            fansinfo.User.Name,
			Avatar:          fansinfo.User.Avatar,
			BackgroundImage: fansinfo.User.BackgroundImage,
			FavoriteCount:   fansinfo.User.FavoriteCount,
			FollowCount:     fansinfo.User.FollowCount,
			FollowerCount:   fansinfo.User.FollowerCount,
			WorkCount:       fansinfo.User.WorkCount,
			Signature:       fansinfo.User.Signature,
			TotalFavorited:  fansinfo.User.TotalFavorited,
			IsFollow:        isFollow,
		})
	}

	return &relation_pb.GetFansListResponse{
		StatusCode: 0, StatusMsg: "粉丝列表",
		UserList: fansMsg,
	}, nil

}
