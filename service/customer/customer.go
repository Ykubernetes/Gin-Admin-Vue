package customer

import (
	"fmt"
	"gitee.com/go-server/global"
	"gitee.com/go-server/models/customer"
	"gitee.com/go-server/models/system"
	"gitee.com/go-server/service/response"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"time"
)

type CustomerListParam struct {
	// 这里跟前端请求一一对应
	Page   int    `form:"page" json:"page"`
	Limit  int    `form:"limit"`
	Sort   string `form:"sort"` // 前端提交的-id 来控制升序降序
	Key    string `form:"key"`  // 对应前端搜索框中的请求
	Status uint8  `form:"status"`
}

// detail 信息返回部分数据
type CustomerResponeDetail struct {
	Id             string `json:"id" form:"id"`
	UserName       string `json:"user_name" form:"user_name"`
	CompanyName    string `json:"company_name" form:"company_name"`
	Position       string `json:"position" form:"position"`
	Birthday       string `json:"birthday" form:"birthday"`
	Gender         string `json:"gender" form:"gender"`
	Province       string `json:"province" form:"province"`
	City           string `json:"city" form:"city"`
	County         string `json:"county" form:"county"`
	Address        string `json:"address" form:"address"`
	Email          string `json:"email" form:"email"`
	Phone          string `json:"phone" form:"phone"`
	Founder        string `json:"founder" form:"founder"`
	CustomerSoruce string `json:"customer_soruce" form:"customer_soruce"`
	Status         uint8  `json:"status" form:"status"`
	Remark         string `json:"remark" form:"remark"`
	CreatedAt      string `json:"create_at" form:"create_at"`
}

// GetCustomerList 客户信息分页数据
func GetCustomerList(c *gin.Context) {
	var requestListParam CustomerListParam
	// 对应前端GET方法
	err := c.ShouldBindQuery(&requestListParam)
	if err != nil {
		response.ResFail(c, "参数绑定错误...")
		return
	}
	// 构造分页条件
	var whereOrder []system.PageOrder
	order := "ID DESC"
	if len(requestListParam.Sort) >= 2 {
		orderType := requestListParam.Sort[0:1]
		order = requestListParam.Sort[1:len(requestListParam.Sort)]
		if orderType == "+" {
			order += " ASC"
		} else {
			order += " DESC"
		}
	}
	whereOrder = append(whereOrder, system.PageOrder{Order: order})
	// Key 有值 代表前端使用搜索框进行模型查询了
	// 支持几个or like查询 就需要增加arr = append(arr, v)
	if requestListParam.Key != "" {
		v := "%" + requestListParam.Key + "%"
		var arr []interface{}
		arr = append(arr, v)
		arr = append(arr, v)
		arr = append(arr, v)
		whereOrder = append(whereOrder, system.PageOrder{Where: "user_name like ? or phone like ? or founder like ?", Value: arr})
	}
	// 对应前端搜索框后面的状态选择框
	if requestListParam.Status > 0 {
		var arr []interface{}
		arr = append(arr, requestListParam.Status)
		whereOrder = append(whereOrder, system.PageOrder{Where: "status = ?", Value: arr})
	}
	// 分页total 对应前端分页总数
	var total int64
	list := []customer.Customer{}
	table := global.DB.Table("customer")
	err = system.GetPage(table, &customer.Customer{}, &list, requestListParam.Page, requestListParam.Limit, &total, whereOrder...)
	if err != nil {
		response.ResErrSrv(c, err)
		return
	}
	response.ResSuccessPage(c, total, &list)
}

// CustomerDetail 客户信息详情 对应前端操作--》查看按钮接口
func CustomerDetail(c *gin.Context) {
	queryId, b := c.GetQuery("id")
	if !b {
		return
	}
	var customerUser customer.Customer
	err := global.DB.Where("id = ?", queryId).First(&customerUser).Error
	if err == gorm.ErrRecordNotFound {
		global.Log.Warn("数据库客户表中未发现查询的id")
		return
	} else if err != nil {
		global.Log.Warnf("数据库客户表中查询发生错误:%s", err)
		response.ResErrSrv(c, err)
		return
	}
	// 返回一部分数据 不要将全部结构体中的数据进行返回
	resData := CustomerResponeDetail{
		Id:             queryId,
		UserName:       customerUser.UserName,
		CompanyName:    customerUser.CompanyName,
		Position:       customerUser.Position,
		Birthday:       customerUser.Birthday,
		Gender:         customerUser.Gender,
		Province:       customerUser.Province,
		City:           customerUser.City,
		County:         customerUser.County,
		Address:        customerUser.Address,
		Email:          customerUser.Email,
		Phone:          customerUser.Phone,
		Founder:        customerUser.Founder,
		CustomerSoruce: customerUser.CustomerSoruce,
		Status:         customerUser.Status,
		Remark:         customerUser.Remark,
		CreatedAt:      customerUser.CreatedAt.Format("2006-01-02 15:04:05"),
	}
	response.ResSuccess(c, &resData)
}

// CustomerAdd 客户信息添加
func CreateCustomer(c *gin.Context) {
	u1 := CustomerResponeDetail{}
	err := c.ShouldBindJSON(&u1)
	if err != nil {
		response.ResErrSrv(c, err)
		return
	}
	birthdayFormart, err := time.Parse("20060102", u1.Birthday)
	if err != nil {
		response.ResFail(c, "出生日期格式不正确")
		return
	}
	dateF := birthdayFormart.Format("2006年1月2日")
	// 在添加前 要确认数据中 不存在重复
	var cus customer.Customer
	err = global.DB.Where("company_name = ?", u1.CompanyName).First(&cus).Error
	if err == gorm.ErrRecordNotFound { // record not found
		global.Log.Warnln("数据库中没有公司名称，record not found")
		var customer = customer.Customer{
			UserName:       u1.UserName,
			CompanyName:    u1.CompanyName,
			Position:       u1.Position,
			Birthday:       fmt.Sprintf("%s", dateF), //
			Gender:         u1.Gender,
			Province:       u1.Province,
			City:           u1.City,
			County:         u1.County,
			Address:        u1.Address,
			Email:          u1.Email,
			Phone:          u1.Phone,
			Founder:        u1.Founder,
			CustomerSoruce: u1.CustomerSoruce,
			Status:         u1.Status,
			Remark:         u1.Remark,
		}
		err = global.DB.Create(&customer).Error
		if err != nil {
			global.Log.Warnf("创建客户信息失败,%s", err)
			response.ResFail(c, "创建客户信息失败")
			return
		}
		response.ResSuccess(c, gin.H{"id": u1.UserName})
	} else if err != nil {
		global.Log.Warnln("err里面有其他错误")
		return
	} else {
		response.ResFail(c, "已存在当前公司")
		return
	}
}

// UpdateCustomer 客户信息更新
func UpdateCustomer(c *gin.Context) {
	u1 := CustomerResponeDetail{}
	err := c.ShouldBindJSON(&u1)
	if err != nil {
		response.ResErrSrv(c, err)
		return
	}
	var customerUser customer.Customer

	// 通过Id 查询出用户
	err = global.DB.Where("id = ?", u1.Id).First(&customerUser).Error
	if err == gorm.ErrRecordNotFound {
		global.Log.Warn("数据库客户表中未发现查询的客户Id")
		return
	} else if err != nil {
		global.Log.Warnf("数据库客户表中查询发生错误:%s", err)
		response.ResErrSrv(c, err)
		return
	}

	// 校验更新传参的数据 这里先滤过
	global.DB.Model(&customerUser).Updates(customer.Customer{
		UserName:       u1.UserName,
		CompanyName:    u1.CompanyName,
		Position:       u1.Position,
		Birthday:       u1.Birthday,
		Gender:         u1.Gender,
		Province:       u1.Province,
		City:           u1.City,
		County:         u1.County,
		Address:        u1.Address,
		Email:          u1.Email,
		Phone:          u1.Phone,
		Founder:        u1.Founder,
		CustomerSoruce: u1.CustomerSoruce,
		Status:         u1.Status,
		Remark:         u1.Remark,
	})
	response.ResSuccessMsg(c)
}

type CustomerDelParam struct {
	CustomerId []uint64 `json:"customerid" binding:"required"`
}

// DeleteCustomer 删除客户
func DeleteCustomer(c *gin.Context) {
	// 获取需要删除的id 有可能是批量删除
	var customerId CustomerDelParam
	if err := c.ShouldBindJSON(&customerId); err != nil {
		response.ResFail(c, "删除客户参数有误...")
		return
	}
	var uList []uint64
	for _, id := range customerId.CustomerId {
		uList = append(uList, id)
	}
	u1 := customer.Customer{}
	err := global.DB.Delete(&u1, uList).Error
	if err != nil {
		response.ResErrSrv(c, err)
		return
	}
	response.ResSuccessMsg(c)

}
