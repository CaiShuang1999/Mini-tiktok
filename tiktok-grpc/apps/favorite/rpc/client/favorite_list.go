package client

import (
	"log"
	"net/http"
	"tiktok-grpc/apps/favorite/favorite_pb"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func FavoriteList(c *gin.Context) {
	serviceHost := viper.GetString("favoriteService.host")
	servicePort := viper.GetString("favoriteService.port")
	serviceAddress := serviceHost + ":" + servicePort
	conn, err := grpc.Dial(serviceAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	rpcclient := favorite_pb.NewFavoriteServiceClient(conn)
	req := &favorite_pb.FavoriteListRequest{

		Token:  c.Query("token"),
		UserId: c.Query("user_id"),
	}
	res, err := rpcclient.FavoriteList(c, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}
