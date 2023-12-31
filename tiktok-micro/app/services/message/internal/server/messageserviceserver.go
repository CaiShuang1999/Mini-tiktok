// Code generated by goctl. DO NOT EDIT.
// Source: messages.proto

package server

import (
	"context"

	"tiktok-micro/app/services/message/internal/logic"
	"tiktok-micro/app/services/message/internal/svc"
	"tiktok-micro/app/services/message/messages_pb"
)

type MessageServiceServer struct {
	svcCtx *svc.ServiceContext
	messages_pb.UnimplementedMessageServiceServer
}

func NewMessageServiceServer(svcCtx *svc.ServiceContext) *MessageServiceServer {
	return &MessageServiceServer{
		svcCtx: svcCtx,
	}
}

func (s *MessageServiceServer) MessageAction(ctx context.Context, in *messages_pb.MessageActionRequest) (*messages_pb.MessageActionResponse, error) {
	l := logic.NewMessageActionLogic(ctx, s.svcCtx)
	return l.MessageAction(in)
}

func (s *MessageServiceServer) GetMessageList(ctx context.Context, in *messages_pb.MessageListRequest) (*messages_pb.MessageListResponse, error) {
	l := logic.NewGetMessageListLogic(ctx, s.svcCtx)
	return l.GetMessageList(in)
}
