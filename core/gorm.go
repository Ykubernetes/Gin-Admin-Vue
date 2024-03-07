package core

import (
	"fmt"
	"gitee.com/go-server/global"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"time"
)

func InitGorm() *gorm.DB {
	return MysqlConnect()
}

func MysqlConnect() *gorm.DB {
	if global.Config.Mysql.Host == "" {
		global.Log.Warn("未配置Mysql，取消Gorm连接...")
		return nil
	}
	dsn := global.Config.Mysql.Dsn()

	var mysqlLogger logger.Interface
	if global.Config.System.Env == "dev" {
		// 开发环境显示所有SQL
		mysqlLogger = logger.Default.LogMode(logger.Info)
	} else {
		mysqlLogger = logger.Default.LogMode(logger.Error) // 只打印错误的sql
	}

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{SingularTable: true}, //使用单数表名
		Logger:         mysqlLogger,
	})
	if err != nil {
		global.Log.Error(fmt.Sprintf("[%s] Mysql数据库连接失败..."))
	}
	sqlDB, _ := db.DB()
	sqlDB.SetMaxIdleConns(10)               //最大空闲连接数
	sqlDB.SetMaxOpenConns(100)              // 最多可容纳
	sqlDB.SetConnMaxLifetime(time.Hour * 4) // 连接最大复用时间，不能超过mysql的wait_timeout
	return db
}
