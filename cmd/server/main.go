package main

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
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

func initDb() {
	dbHost := viper.GetString("DB_HOST")
	dbPort := viper.GetString("DB_PORT")
	dbUser := viper.GetString("DB_USERNAME")
	dbPass := viper.GetString("DB_PASSWORD")
	dbCharset := viper.GetString("DB_CHARSET")
	dbLocal := viper.GetString("DB_LOCAL")
	dbDatabase := viper.GetString("DB_DATABASE")
	fmt.Println(dbHost, dbPort, dbUser, dbPass, dbCharset, dbLocal, dbDatabase)
	dbDsn := dbUser + ":" +
		dbPass + "@tcp(" + dbHost + ":" + dbPort + ")/" + dbDatabase +
		"?charset=" + dbCharset +
		"&parseTime=True&loc=" + dbLocal
	_, err := gorm.Open(mysql.New(mysql.Config{
		DSN:                      dbDsn, // DSN data source name
		DefaultStringSize:        256,   // string 类型字段的默认长度
		DisableDatetimePrecision: true,  // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
		DontSupportRenameIndex:   true,  // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
		DontSupportRenameColumn:  true,  // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
	}), &gorm.Config{})
	if err != nil { // Handle errors reading the config file
		log.Error("init db failed err: ", err)
		panic(err)
	}
}

func init() {
	initLogger()
	initViper()
	initDb()
}

func main() {
	gin.SetMode(viper.GetString("RUN_MODE"))
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.Run("127.0.0.1:8080") // 监听并在 0.0.0.0:8080 上启动服务
}
