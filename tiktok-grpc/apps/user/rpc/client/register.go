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

func Register(c *gin.Context) {
	serviceHost := viper.GetString("userService.host")
	servicePort := viper.GetString("userService.port")
	serviceAddress := serviceHost + ":" + servicePort
	conn, err := grpc.Dial(serviceAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	rpcclient := pb.NewUserServiceClient(conn)
	req := &pb.RegisterRequest{
		Username: c.Query("username"),
		Password: c.Query("password"),
	}
	res, err := rpcclient.Register(c, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Message传回gin的JSON
	c.JSON(http.StatusOK, res)

}
