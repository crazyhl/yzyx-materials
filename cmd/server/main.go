package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/asaskevich/EventBus"
	"github.com/crazyhl/yzyx-materials/internal/bus"
	"github.com/crazyhl/yzyx-materials/internal/db"
	_ "github.com/crazyhl/yzyx-materials/internal/validator"
	"github.com/crazyhl/yzyx-materials/module/account"
	"github.com/crazyhl/yzyx-materials/module/domain/models"
	"github.com/crazyhl/yzyx-materials/route"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/gzip"
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
	viper.AutomaticEnv()
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		log.Error("read config file err: ", err)
		panic(err)
	}
}

func init() {
	initLogger()
	initViper()
	db.InitDb()
	models.AutoMigrate()
}

func main() {
	bus.Bus = EventBus.New()
	bus.Bus.Subscribe("account:updateProfit", account.UpdateAccountProfit)
	// init gin
	gin.SetMode(viper.GetString("RUN_MODE"))
	router := gin.Default()
	router.Use(gzip.Gzip(gzip.DefaultCompression))
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{viper.GetString("CORS_HOST")},
		AllowMethods:     []string{"PUT", "GET", "POST", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Authorization", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	route.InitRouter(router)

	srv := &http.Server{
		Addr:         ":" + viper.GetString("SERVER_PORT"),
		Handler:      router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
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
	bus.Bus.Unsubscribe("account:updateProfit", account.UpdateAccountProfit)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	log.Println("Server exiting")

}
