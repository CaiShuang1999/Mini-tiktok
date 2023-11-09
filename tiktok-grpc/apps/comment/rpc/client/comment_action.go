package client

import (
	"log"
	"net/http"
	"tiktok-grpc/apps/comment/comment_pb"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func CommentAction(c *gin.Context) {
	serviceHost := viper.GetString("commentService.host")
	servicePort := viper.GetString("commentService.port")
	serviceAddress := serviceHost + ":" + servicePort

	conn, err := grpc.Dial(serviceAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	rpcclient := comment_pb.NewCommentServiceClient(conn)
	req := &comment_pb.CommentActionRequest{
		VideoId:     c.Query("video_id"),
		Token:       c.Query("token"),
		ActionType:  c.Query("action_type"),
		CommentText: c.Query("comment_text"),
		CommentId:   c.Query("comment_id"),
	}
	res, err := rpcclient.CommentAction(c, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, res)

}
