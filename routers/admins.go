package routers

import (
	"gitee.com/go-server/service/system"
	"github.com/gin-gonic/gin"
)

func AdminsRouter(r *gin.Engine) {
	r.GET("/admins/list", system.GetAdminsList)
	r.GET("/admins/adminsroleidlist", system.GetAdminsRoleIDList)
	r.GET("/admins/detail", system.GetAdminsDetail)
	r.PUT("/admins/update", system.AdminsUpdate)
	r.POST("/admins/create", system.AdminsCreate)
	r.DELETE("/admins/delete", system.AdminsDelete)
	r.POST("/admins/setrole", system.AdminsSetrole)
}
