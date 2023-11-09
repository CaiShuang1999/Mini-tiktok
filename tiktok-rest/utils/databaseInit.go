package utils

import (
	"context"
	"fmt"
	"tiktok-rest/model"

	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// viper配置数据库config初始化
func InitConfig() {

	viper.SetConfigName("database") //config文件名
	viper.SetConfigType("yaml")     //config文件类型（yaml，json）及其他类型
	viper.AddConfigPath("config")   //config路径

	//读取Config
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("成功读取配置文件!")

}

var DB *gorm.DB //全局变量 使用DB将数据库传给外部文件

func InitMysql() {
	//获取Config
	dbHost := viper.GetString("database.host")
	dbPort := viper.GetInt("database.port")
	dbName := viper.GetString("database.dbname")
	dbUsername := viper.GetString("database.username")
	dbPassword := viper.GetString("database.password")
	dbCharset := viper.GetString("database.charset")
	dbParseTime := viper.GetBool("database.parseTime")
	dbLoc := viper.GetString("database.loc")

	// 构建数据库连接字符串
	dbConnectionString := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=%t&loc=%s", dbUsername, dbPassword, dbHost, dbPort, dbName, dbCharset, dbParseTime, dbLoc)

	// 在这里使用数据库连接字符串连接数据库
	fmt.Println("Database Connection String:", dbConnectionString)

	db, err := gorm.Open(mysql.Open(dbConnectionString), &gorm.Config{})
	if err != nil {
		fmt.Println("数据库连接失败:", err)
		return
	}
	fmt.Println("数据库连接成功！")

	err = db.AutoMigrate(&model.User{}, &model.Video{}, &model.Favorite{},
		&model.Comment{}, &model.Relation{}, &model.Message{})
	if err != nil {
		fmt.Println("数据表迁移错误")
		fmt.Println(err) // 错误处理
		return
	}
	fmt.Println("数据表迁移成功")
	//传递数据库
	DB = db

}

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
