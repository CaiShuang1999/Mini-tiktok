package logic

import (
	"context"
	"strconv"
	"time"

	"tiktok-micro/app/services/message/internal/svc"
	"tiktok-micro/app/services/message/messages_pb"
	"tiktok-micro/common/jwtx"
	"tiktok-micro/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type MessageActionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewMessageActionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MessageActionLogic {
	return &MessageActionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *MessageActionLogic) MessageAction(in *messages_pb.MessageActionRequest) (*messages_pb.MessageActionResponse, error) {
	token := in.Token
	toUserIDstr := in.ToUserId
	actionType := in.ActionType
	content := in.Content

	if actionType != "1" {
		return &messages_pb.MessageActionResponse{
			StatusCode: 1, StatusMsg: "无效actionType",
		}, nil

	}

	tokenMsg, ok := jwtx.ParseToken(token)
	if !ok {
		return &messages_pb.MessageActionResponse{
			StatusCode: 1, StatusMsg: "无效token",
		}, nil
	}

	toUserID, err := strconv.Atoi(toUserIDstr)
	if err != nil {
		return &messages_pb.MessageActionResponse{
			StatusCode: 1, StatusMsg: "无效ID",
		}, nil
	}

	var newMessage model.Message
	newMessage.FromUserID = tokenMsg.UserID
	newMessage.ToUserID = int64(toUserID)
	newMessage.Content = content
	newMessage.CreateTime = time.Now().Unix()

	db := l.svcCtx.DB
	db.Create(&newMessage)

	return &messages_pb.MessageActionResponse{
		StatusCode: 0, StatusMsg: "发送信息成功",
	}, nil
}
