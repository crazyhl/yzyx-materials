package main

import (
	"os"

	"github.com/crazyhl/yzyx-materials/internal"
	"github.com/crazyhl/yzyx-materials/module/user"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var runDir string

// initLogger 初始化日志配置
func initLogger() {
	// 设置日志输出格式
	log.SetFormatter(&log.JSONFormatter{})
	// 获取程序运行目录
	var err error
	runDir, err = os.Getwd()
	if err != nil {
		panic(err)
	}
	// 创建日志目录
	err = os.MkdirAll(runDir+"/log", os.ModePerm)
	if err != nil {
		panic(err)
	}
	// 设置日志文件输出路径
	file, err := os.OpenFile(runDir+"/log/runLog.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		log.SetOutput(file)
	} else {
		panic(err)
	}

}

func initViper() {
	// 配置 viper 信息，并读入配置
	viper.SetConfigName("env")
	viper.SetConfigType("env")
	viper.AddConfigPath(runDir)
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		log.Error("read config file err: ", err)
		panic(err)
	}
}

func init() {
	initLogger()
	initViper()
	internal.InitDb()
	user.AutoMigrate()
}

func main() {
	gin.SetMode(viper.GetString("RUN_MODE"))
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.POST("/user/register", user.Register)
	r.Run("127.0.0.1:8080") // 监听并在 0.0.0.0:8080 上启动服务
}
