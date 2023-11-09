package client

import (
	"log"
	"net/http"
	"tiktok-grpc/apps/user/pb"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func GetUserInfo(c *gin.Context) {
	serviceHost := viper.GetString("userService.host")
	servicePort := viper.GetString("userService.port")
	serviceAddress := serviceHost + ":" + servicePort

	conn, err := grpc.Dial(serviceAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	rpcclient := pb.NewUserServiceClient(conn)
	req := &pb.UserInfoRequest{
		UserId: c.Query("user_id"),
		Token:  c.Query("token"),
	}
	res, err := rpcclient.GetUser(c, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, res)

}
