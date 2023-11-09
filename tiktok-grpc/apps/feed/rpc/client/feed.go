package client

import (
	"log"
	"net/http"
	"tiktok-grpc/apps/feed/feed_pb"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func FeedList(c *gin.Context) {
	serviceHost := viper.GetString("feedService.host")
	servicePort := viper.GetString("feedService.port")
	serviceAddress := serviceHost + ":" + servicePort
	conn, err := grpc.Dial(serviceAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	rpcclient := feed_pb.NewFeedServiceClient(conn)
	req := &feed_pb.FeedRequest{
		LatestTime: c.Query("latest_time"),
		Token:      c.Query("token"),
	}
	res, err := rpcclient.FeedList(c, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)

}
