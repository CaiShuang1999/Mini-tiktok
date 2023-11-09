package apps

import (
	"fmt"
	comment_server "tiktok-grpc/apps/comment/rpc/server"
	favorite_server "tiktok-grpc/apps/favorite/rpc/server"
	feed_server "tiktok-grpc/apps/feed/rpc/server"
	message_server "tiktok-grpc/apps/message/rpc/server"
	relation_server "tiktok-grpc/apps/relation/rpc/server"
	user_server "tiktok-grpc/apps/user/rpc/server"
	video_server "tiktok-grpc/apps/video/rpc/server"

	"github.com/spf13/viper"
)

func ServerInits() {
	viper.SetConfigName("services") //config文件名
	viper.SetConfigType("yaml")     //config文件类型（yaml，json）及其他类型
	viper.AddConfigPath("configs")  //config路径

	//读取Config
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("成功读取微服务配置文件!")

	go user_server.UserServerInit() //用户微服务服务器启动

	go video_server.VideoServerInit() //视频微服务服务器启动

	go feed_server.FeedServerInit() //feed

	go favorite_server.FavoriteServerInit() //点赞

	go comment_server.CommentServerInit() //评论

	go relation_server.RelationServerInit() //关系

	go message_server.MessageServerInit() //聊天
}
