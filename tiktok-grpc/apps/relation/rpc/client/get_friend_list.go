package client

import (
	"log"
	"net/http"
	"tiktok-grpc/apps/relation/relation_pb"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func GetFriendList(c *gin.Context) {
	serviceHost := viper.GetString("relationService.host")
	servicePort := viper.GetString("relationService.port")
	serviceAddress := serviceHost + ":" + servicePort

	conn, err := grpc.Dial(serviceAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	rpcclient := relation_pb.NewRelationServiceClient(conn)
	req := &relation_pb.GetFriendListRequest{
		UserId: c.Query("user_id"),
		Token:  c.Query("token"),
	}
	res, err := rpcclient.GetFriendList(c, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, res)

}
