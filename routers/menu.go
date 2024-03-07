package routers

import (
	"gitee.com/go-server/service/system" // handler层 不是model层的system
	"github.com/gin-gonic/gin"
)

func MenuRouter(r *gin.Engine) {
	r.GET("/menu/allmenu", system.GetAllMenu)
	r.GET("/menu/list", system.GetMenuList)
	r.GET("/menu/detail", system.Detail)
	r.GET("/menu/menubuttonlist", system.MenuButtonList)
	r.POST("/menu/create", system.Create)
	r.DELETE("/menu/delete", system.Delete)
	r.PUT("/menu/update", system.Update)
}
