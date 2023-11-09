package server

import (
	"context"
	"strconv"
	"tiktok-grpc/apps/message/messages_pb"
	"tiktok-grpc/cmd"
	"tiktok-grpc/common/jwtx"
	"tiktok-grpc/common/utils"
	"tiktok-grpc/model"
	"time"
)

type MessageServiceServer struct {
	messages_pb.UnimplementedMessageServiceServer
}

func (p *MessageServiceServer) MessageAction(ctx context.Context, req *messages_pb.MessageActionRequest) (*messages_pb.MessageActionResponse, error) {
	token := req.Token
	toUserIDstr := req.ToUserId
	actionType := req.ActionType
	content := req.Content

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

	db := cmd.DB
	db.Create(&newMessage)

	return &messages_pb.MessageActionResponse{
		StatusCode: 0, StatusMsg: "发送信息成功",
	}, nil

}
func (p *MessageServiceServer) GetMessageList(ctx context.Context, req *messages_pb.MessageListRequest) (*messages_pb.MessageListResponse, error) {
	token := req.Token
	toUserIDstr := req.ToUserId
	preMsgTimeStr := req.PreMsgTime
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
	db := cmd.DB
	db.Where("(from_user_id=? AND to_user_id=? AND create_time > ? ) OR (from_user_id=? AND to_user_id=? AND create_time > ? )",
		tokenMsg.UserID, toUserID, preMsgTime, toUserID, tokenMsg.UserID, preMsgTime).Find(&messageList)

	messageListMsg := []*messages_pb.MessageInfo{}

	for i := range messageList {
		messageListMsg = append(messageListMsg, utils.ConvertMessageToMessageProto(messageList[i]))
	}

	return &messages_pb.MessageListResponse{
		StatusCode: 0, StatusMsg: "聊天列表",
		MessageList: messageListMsg,
	}, nil
}
