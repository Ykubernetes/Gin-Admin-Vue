package routers

import (
	"gitee.com/go-server/service/system"
	"github.com/gin-gonic/gin"
)

func DepartMentRouter(r *gin.Engine) {
	// 部门列表
	r.GET("/department/list", system.GetDepartMentList)
	// 根据id获取部门信息
	r.GET("/department/:id", system.GetDept)
	// 添加部门
	r.POST("/department/create", system.CreateDepartMent)
	// 修改部门
	r.PUT("/department/update", system.UpdateDepartMent)
	// 删除部门
	r.DELETE("/department/del/:id", system.DeleteDepartMent)
	// 查询部门下拉树结构
	r.GET("/department/deptTree", system.GetDeptTree)
}
