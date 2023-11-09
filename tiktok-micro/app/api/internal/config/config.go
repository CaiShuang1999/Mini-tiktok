package config

import (
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	rest.RestConf
	UserRpc     zrpc.RpcClientConf
	VideoRpc    zrpc.RpcClientConf
	FeedRpc     zrpc.RpcClientConf
	FavoriteRpc zrpc.RpcClientConf
	CommentRpc  zrpc.RpcClientConf
	RelationRpc zrpc.RpcClientConf
	MessageRpc  zrpc.RpcClientConf
}
