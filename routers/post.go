package routers

import (
	"gitee.com/go-server/service/system"
	"github.com/gin-gonic/gin"
)

func PostRouter(r *gin.Engine) {
	// 岗位列表
	r.GET("/post/list", system.GetPostList)
	// 根据id获取岗位信息
	r.GET("/post/:id", system.GetPost)
	// 添加岗位
	r.POST("/post/create", system.CreatePost)
	// 修改岗位
	r.PUT("/post/update", system.UpdatePost)
	// 删除部门
	r.DELETE("/post/del/:id", system.DeletePost)
}
