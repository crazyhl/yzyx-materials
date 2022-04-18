package internal

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDb() {
	dbHost := viper.GetString("DB_HOST")
	dbPort := viper.GetString("DB_PORT")
	dbUser := viper.GetString("DB_USERNAME")
	dbPass := viper.GetString("DB_PASSWORD")
	dbCharset := viper.GetString("DB_CHARSET")
	dbLocal := viper.GetString("DB_LOCAL")
	dbDatabase := viper.GetString("DB_DATABASE")
	dbDsn := dbUser + ":" +
		dbPass + "@tcp(" + dbHost + ":" + dbPort + ")/" + dbDatabase +
		"?charset=" + dbCharset +
		"&parseTime=True&loc=" + dbLocal
	var err error
	DB, err = gorm.Open(mysql.New(mysql.Config{
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