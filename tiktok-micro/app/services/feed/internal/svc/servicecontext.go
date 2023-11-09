package svc

import (
	"context"
	"fmt"
	"tiktok-micro/app/services/feed/internal/config"

	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type ServiceContext struct {
	Config config.Config
	DB     *gorm.DB
	Rdb_0  *redis.Client
	Rdb_1  *redis.Client
}

func NewServiceContext(c config.Config) *ServiceContext {
	//mysql连接
	db, err := MysqlDBinit(c.MysqlDB.DataSource)
	if err != nil {
		fmt.Println("mysql连接错误:" + err.Error())
		panic(err)
	}

	//redis连接
	rdb0, err := RedisDBinit(c.RedisDB_0.Addr, c.RedisDB_0.DB)
	if err != nil {
		fmt.Println("redis1连接错误:" + err.Error())
		panic(err)
	}

	rdb1, err := RedisDBinit(c.RedisDB_1.Addr, c.RedisDB_1.DB)
	if err != nil {
		fmt.Println("redis2连接错误:" + err.Error())
		panic(err)
	}
	return &ServiceContext{
		Config: c,
		DB:     db,
		Rdb_0:  rdb0,
		Rdb_1:  rdb1,
	}
}
func MysqlDBinit(dbConnectionString string) (*gorm.DB, error) {
	db, err := gorm.Open(mysql.Open(dbConnectionString), &gorm.Config{})
	if err != nil {
		fmt.Println("数据库连接失败!")
		return nil, err
	}
	fmt.Println("数据库连接成功！")
	/*
		err = db.AutoMigrate(&model.User{}, &model.Video{}, &model.Favorite{},
			&model.Comment{}, &model.Relation{}, &model.Message{})
		if err != nil {
			fmt.Println("数据表迁移错误!")
			return nil, err
		}
		fmt.Println("数据表迁移成功!")
	*/
	return db, nil
}

func RedisDBinit(addr string, db int) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr: addr, // Redis 服务器地址
		DB:   db,   // 使用默认的数据库
	})
	ctx := context.Background()
	// 使用 Ping 方法检查与 Redis 服务器的连接是否成功
	_, err := client.Ping(ctx).Result()
	if err != nil {

		return nil, err
	}

	fmt.Printf("Redis: %d, 连接成功!\n", db)

	return client, nil
}
