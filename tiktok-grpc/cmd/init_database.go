package cmd

import (
	"fmt"
	"tiktok-grpc/model"

	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB //全局变量 使用DB将数据库传给外部文件
func InitSQL() {
	viper.SetConfigName("database") //config文件名
	viper.SetConfigType("yaml")     //config文件类型（yaml，json）及其他类型
	viper.AddConfigPath("configs")  //config路径

	//读取Config
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("成功读取sql配置文件!")

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
