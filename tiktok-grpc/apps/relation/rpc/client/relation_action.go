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

func RelationAction(c *gin.Context) {
	serviceHost := viper.GetString("relationService.host")
	servicePort := viper.GetString("relationService.port")
	serviceAddress := serviceHost + ":" + servicePort

	conn, err := grpc.Dial(serviceAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	rpcclient := relation_pb.NewRelationServiceClient(conn)
	req := &relation_pb.RelationActionRequest{
		Token:      c.Query("token"),
		ToUserId:   c.Query("to_user_id"),
		ActionType: c.Query("action_type"),
	}
	res, err := rpcclient.RelationAction(c, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, res)

}
