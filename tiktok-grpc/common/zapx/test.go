package main

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Logger *zap.Logger

func InitLogger() *zap.Logger {
	file, _ := os.Create("test.log") //文件
	writeSyncer := zapcore.AddSync(file)
	encoder := zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
	core := zapcore.NewCore(encoder, writeSyncer, zapcore.ErrorLevel)
	logger := zap.New(core)
	return logger
}

func main() {

}
