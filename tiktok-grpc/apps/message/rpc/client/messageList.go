package client

import (
	"log"
	"net/http"
	"tiktok-grpc/apps/message/messages_pb"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func GetMessageList(c *gin.Context) {
	serviceHost := viper.GetString("messageService.host")
	servicePort := viper.GetString("messageService.port")
	serviceAddress := serviceHost + ":" + servicePort

	conn, err := grpc.Dial(serviceAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	rpcclient := messages_pb.NewMessageServiceClient(conn)
	req := &messages_pb.MessageListRequest{
		ToUserId:   c.Query("to_user_id"),
		Token:      c.Query("token"),
		PreMsgTime: c.Query("pre_msg_time"),
	}
	res, err := rpcclient.GetMessageList(c, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, res)

}
