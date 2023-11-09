package config

import "github.com/zeromicro/go-zero/zrpc"

type Config struct {
	zrpc.RpcServerConf
	MysqlDB struct {
		DataSource string
	}
	Nginx struct {
		Addr string
	}

	RedisDB_0 struct {
		Addr string
		DB   int
	}
	RedisDB_1 struct {
		Addr string
		DB   int
	}
	VideoRpc zrpc.RpcClientConf
}
