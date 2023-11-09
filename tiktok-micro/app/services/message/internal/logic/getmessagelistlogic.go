package logic

import (
	"context"
	"strconv"

	"tiktok-micro/app/services/message/internal/svc"
	"tiktok-micro/app/services/message/messages_pb"
	"tiktok-micro/common/jwtx"
	"tiktok-micro/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetMessageListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetMessageListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetMessageListLogic {
	return &GetMessageListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetMessageListLogic) GetMessageList(in *messages_pb.MessageListRequest) (*messages_pb.MessageListResponse, error) {
	token := in.Token
	toUserIDstr := in.ToUserId
	preMsgTimeStr := in.PreMsgTime
	preMsgTime, _ := strconv.Atoi(preMsgTimeStr)

	tokenMsg, ok := jwtx.ParseToken(token)
	if !ok {

		return &messages_pb.MessageListResponse{
			StatusCode: 1, StatusMsg: "无效token",
		}, nil
	}

	toUserID, err := strconv.Atoi(toUserIDstr)
	if err != nil {
		return &messages_pb.MessageListResponse{
			StatusCode: 1, StatusMsg: "无效ID",
		}, nil
	}

	var messageList []model.Message
	db := l.svcCtx.DB
	db.Where("(from_user_id=? AND to_user_id=? AND create_time > ? ) OR (from_user_id=? AND to_user_id=? AND create_time > ? )",
		tokenMsg.UserID, toUserID, preMsgTime, toUserID, tokenMsg.UserID, preMsgTime).Find(&messageList)

	messageListMsg := []*messages_pb.MessageInfo{}

	for i := range messageList {
		messageListMsg = append(messageListMsg, &messages_pb.MessageInfo{
			Id:         messageList[i].ID,
			ToUserId:   messageList[i].ToUserID,
			FromUserId: messageList[i].FromUserID,
			CreateTime: messageList[i].CreateTime,
			Content:    messageList[i].Content,
		})
	}

	return &messages_pb.MessageListResponse{
		StatusCode: 0, StatusMsg: "聊天列表",
		MessageList: messageListMsg,
	}, nil
}
