package routers

import (
	"gitee.com/go-server/global"
	"gitee.com/go-server/middleware"
	"github.com/gin-gonic/gin"
)

func InitRoter() *gin.Engine {
	gin.SetMode(global.Config.System.Env)
	r := gin.Default()
	r.Static("/upload/avator_img", "upload/avator_img")
	// cors
	r.Use(middleware.Cors())
	r.Use(middleware.UserAuthMiddleware())
	r.Use(middleware.CasbinCheckMiddleware())
	// 加载静态图片 供前端使用
	SettingsRouter(r)
	LoginRouter(r)
	MenuRouter(r)
	AdminsRouter(r)
	RolesRouter(r)
	DepartMentRouter(r)
	PostRouter(r)
	CustomerRouter(r)
	return r
}
