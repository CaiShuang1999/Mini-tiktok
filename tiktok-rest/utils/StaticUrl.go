package utils

import (
	"strconv"

	"github.com/spf13/viper"
)

var StaticUrl string

func GetStaticUrl() {
	Host := viper.GetString("NginxStatic.host")
	PortInt := viper.GetInt("NginxStatic.port")
	Port := strconv.Itoa(PortInt)
	StaticUrl = "http://" + Host + ":" + Port
}
