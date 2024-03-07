package main

import (
	"context"
	"fmt"
	"gitee.com/go-server/core"
	"gitee.com/go-server/global"
	"gitee.com/go-server/models/customer"
	"gitee.com/go-server/models/system"
	"gitee.com/go-server/routers"
	"net/http"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	// 加载配置文件
	core.InitConf()
	// 初始化日志
	global.Log = core.InitLogger()
	// 数据库连接
	global.DB = core.InitGorm()
	// Redis 连接
	global.RedisDB = core.RedisConn()

	// Migratioon
	global.DB.AutoMigrate(&system.Menu{})
	global.DB.AutoMigrate(&system.Admins{})
	global.DB.AutoMigrate(&system.Dept{})
	global.DB.AutoMigrate(&system.RoleMenu{})
	global.DB.AutoMigrate(&system.Role{})
	global.DB.AutoMigrate(&system.AdminsRole{})
	global.DB.AutoMigrate(&system.Post{})
	global.DB.AutoMigrate(&customer.Customer{})

	// Casbin 初始化 传入数据库中
	core.InitCasbinEnforcer(global.Config.Mysql.Dsn())

	// 创建侦听来自操作系统的中断信号的上下文
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	router := routers.InitRoter()

	srv := &http.Server{
		Addr:    ":" + fmt.Sprintf("%v", global.Config.System.Port),
		Handler: router,
	}

	// 在goroutine中初始化服务器，以便
	//它不会阻止下面优雅的关闭处理
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			global.Log.Fatalf("Server listen: %s\n", err)
		}
	}()
	// 打印程序启动的host和port
	global.Log.Infof("Server Starting on %s", global.Config.System.Addr())
	// 监听终端信号
	<-ctx.Done()
	// 恢复中断信号的默认行为，并通知用户关闭
	stop()
	global.Log.Info("shutting down gracefully, press Ctrl+C again to force")

	//创建延时通知上下文用于通知服务器它还有10秒的时间完成
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		global.Log.Fatal("Server forced to shutdown: ", err)
	} else {
		ExitDelToken2Redis(ctx)
		global.Log.Info("Server exiting")
	}

}

func ExitDelToken2Redis(ctx context.Context) {
	keys, err := global.RedisDB.Keys(ctx, "userIdKey-*").Result()
	if err != nil {
		global.Log.Warnln("退出系统获取用户userIdKey错误:", err)
	}
	for _, key := range keys {
		err := global.RedisDB.Del(ctx, key).Err()
		if err != nil {
			global.Log.Warnf("Failed to delete key %s: %v", key, err)
		} else {
			global.Log.Warnf("Deleted key: %s\n", key)
		}
	}
}
