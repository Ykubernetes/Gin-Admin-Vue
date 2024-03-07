package routers

import (
	"gitee.com/go-server/service/system"
	"github.com/gin-gonic/gin"
)

func RolesRouter(r *gin.Engine) {
	r.GET("/role/allrole", system.GetAllRole)
	r.GET("/role/list", system.GetRoleList)
	r.GET("/role/detail", system.GetRoleDetail)
	r.GET("/role/rolemenuidlist", system.GetRoleMenuIDList)
	r.POST("/role/setrole", system.RoleSetRole)
	r.DELETE("/role/delete", system.RoleDelete)
	r.PUT("/role/update", system.RoleUpdate)
	r.POST("/role/create", system.RoleCreate)
}
