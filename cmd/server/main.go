package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

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
	router := gin.Default()
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	router.POST("/user/register", user.Register)
	router.POST("/user/login", user.Login)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	go func() {
		// 服务连接
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// 等待中断信号以优雅地关闭服务器（设置 5 秒的超时时间）
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	log.Println("Server exiting")

}
