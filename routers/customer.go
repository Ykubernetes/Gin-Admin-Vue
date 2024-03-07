package routers

import (
	"gitee.com/go-server/service/customer"
	"github.com/gin-gonic/gin"
)

func CustomerRouter(r *gin.Engine) {
	// 获取客户列表
	r.GET("/customer/list", customer.GetCustomerList)
	// 获取客户详细信息
	r.GET("/customer/detail", customer.CustomerDetail)
	// 客户信息添加
	r.POST("/customer/create", customer.CreateCustomer)
	// 客户信息更新
	r.PUT("/customer/update", customer.UpdateCustomer)
	// 客户信息删除
	r.DELETE("/customer/delete", customer.DeleteCustomer)
}
