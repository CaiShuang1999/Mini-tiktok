package cmd

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client
var RedisClient2 *redis.Client

func InitRedis() {

	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // Redis 服务器地址
		Password: "",               // 如果有密码，请提供密码
		DB:       0,                // 使用默认的数据库
	})
	ctx := context.Background()
	// 使用 Ping 方法检查与 Redis 服务器的连接是否成功
	_, err := client.Ping(ctx).Result()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	RedisClient = client
	fmt.Println("Redis_01连接成功!")

	client2 := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // Redis 服务器地址
		Password: "",               // 如果有密码，请提供密码
		DB:       1,                // 使用默认的数据库
	})

	// 使用 Ping 方法检查与 Redis 服务器的连接是否成功
	_, err = client2.Ping(ctx).Result()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	RedisClient2 = client2
	fmt.Println("Redis_02连接成功!")

}
