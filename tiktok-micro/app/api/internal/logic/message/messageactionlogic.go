package message

import (
	"context"

	"tiktok-micro/app/api/internal/svc"
	"tiktok-micro/app/api/internal/types"
	"tiktok-micro/app/services/message/messages_pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type MessageActionLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewMessageActionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MessageActionLogic {
	return &MessageActionLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MessageActionLogic) MessageAction(req *types.MessageActionRequest) (resp *messages_pb.MessageActionResponse, err error) {
	res, err := l.svcCtx.MessageRpc.MessageAction(l.ctx, &messages_pb.MessageActionRequest{
		ToUserId:   req.TOUserId,
		Token:      req.Token,
		ActionType: req.ActionType,
		Content:    req.Content,
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}
