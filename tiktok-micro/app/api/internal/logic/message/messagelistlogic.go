package message

import (
	"context"

	"tiktok-micro/app/api/internal/svc"
	"tiktok-micro/app/api/internal/types"
	"tiktok-micro/app/services/message/messages_pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type MessageListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewMessageListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MessageListLogic {
	return &MessageListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MessageListLogic) MessageList(req *types.MessageListRequest) (resp *messages_pb.MessageListResponse, err error) {
	res, err := l.svcCtx.MessageRpc.GetMessageList(l.ctx, &messages_pb.MessageListRequest{
		ToUserId:   req.ToUserId,
		Token:      req.Token,
		PreMsgTime: req.PreMsgTime,
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}
