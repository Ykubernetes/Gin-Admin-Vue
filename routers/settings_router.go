package routers

import (
	"gitee.com/go-server/api"
	"github.com/gin-gonic/gin"
)

func SettingsRouter(r *gin.Engine) {
	settingsApi := api.ApiGroupApp.SettingsApi
	r.GET("/ping", settingsApi.SettingsInfoView)
}
