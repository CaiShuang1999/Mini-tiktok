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

func MessageAction(c *gin.Context) {
	serviceHost := viper.GetString("messageService.host")
	servicePort := viper.GetString("messageService.port")
	serviceAddress := serviceHost + ":" + servicePort

	conn, err := grpc.Dial(serviceAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	rpcclient := messages_pb.NewMessageServiceClient(conn)
	req := &messages_pb.MessageActionRequest{
		ToUserId:   c.Query("to_user_id"),
		Token:      c.Query("token"),
		ActionType: c.Query("action_type"),
		Content:    c.Query("content"),
	}
	res, err := rpcclient.MessageAction(c, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, res)

}
