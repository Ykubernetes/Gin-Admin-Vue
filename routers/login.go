package routers

import (
	"gitee.com/go-server/service/system"
	"github.com/gin-gonic/gin"
)

func LoginRouter(r *gin.Engine) {
	r.POST("/user/login", system.Login)
	r.POST("/user/logout", system.Logout)
	r.POST("/user/register", system.Register)
	r.GET("/user/info", system.GetUserInfo)
	r.GET("/user/captcha/image", system.Getcaptcha)
	r.PUT("/user/editpwd", system.ChangePassword)
	r.PUT("/user/editinfo", system.ChangeUserinfo)
	r.GET("/user/systeminfo", system.GetSystemInfo)
	r.GET("/user/systemstate", system.SystemState)
	// 文件上传-单文件上传
	r.POST("/user/avatar/img", system.AvatorImgUpload)
}
