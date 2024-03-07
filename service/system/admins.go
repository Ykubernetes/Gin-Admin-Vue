package system

import (
	"fmt"
	"gitee.com/go-server/core"
	"gitee.com/go-server/global"
	"gitee.com/go-server/models/system"
	"gitee.com/go-server/service/response"
	"gitee.com/go-server/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"strconv"
)

type AdminsListParam struct {
	Page   int    `form:"page" json:"page"`
	Limit  int    `form:"limit" json:"limit"`
	Sort   string `form:"sort" json:"sort"`
	Key    string `form:"key" json:"key"`
	Status uint8  `form:"status" json:"status"`
}

type AdminsDelParam struct {
	Id []uint64 `json:"id" binding:"required"`
}

type RoleIdRequestParam struct {
	RoleId []uint64 `json:"roleid" binding:"required"`
}

// AdminsResponeList 要返回的用户列表结构体
type AdminsResponeList struct {
	system.Admins        // 要返回的用户列表
	DeptName      string `json:"deptName" gorm:"dept_name"` // 外加用户表和部门表join之后需要的部门名称字段
}

// 分页数据
func GetAdminsList(c *gin.Context) {
	var adminsListParam AdminsListParam
	err := c.ShouldBindQuery(&adminsListParam)
	if err != nil {
		response.ResFail(c, "参数绑定错误...")
		return
	}
	var whereOrder []system.PageOrder
	order := "ID DESC"
	if len(adminsListParam.Sort) >= 2 {
		orderType := adminsListParam.Sort[0:1]
		order = adminsListParam.Sort[1:len(adminsListParam.Sort)]
		if orderType == "+" {
			order += " ASC"
		} else {
			order += " DESC"
		}
	}
	whereOrder = append(whereOrder, system.PageOrder{Order: order})
	if adminsListParam.Key != "" {
		v := "%" + adminsListParam.Key + "%"
		var arr []interface{}
		arr = append(arr, v)
		arr = append(arr, v)
		whereOrder = append(whereOrder, system.PageOrder{Where: "user_name like ? or real_name like ?", Value: arr})
	}
	if adminsListParam.Status > 0 {
		var arr []interface{}
		arr = append(arr, adminsListParam.Status)
		whereOrder = append(whereOrder, system.PageOrder{Where: "status = ?", Value: arr})
	}
	var total int64
	list := []AdminsResponeList{} // 返回的结构体
	// 两张表join之后再做where
	table := global.DB.Select("admins.*,dept_name").Table("admins").Joins("LEFT JOIN department on department.id = admins.dept_id")
	err = system.GetPage(table, &system.Admins{}, &list, adminsListParam.Page, adminsListParam.Limit, &total, whereOrder...)
	if err != nil {
		response.ResErrSrv(c, err)
		return
	}
	// 这里暴露了密码 返回前台数据时,所以需要将密码隐藏
	list1 := []AdminsResponeList{}
	for _, val := range list {
		val.Password = "************"
		list1 = append(list1, val)
	}

	response.ResSuccessPage(c, total, &list1)
}

// 详情
func GetAdminsDetail(c *gin.Context) {
	query, b := c.GetQuery("id")
	if !b {
		return
	}
	var admins system.Admins
	err := global.DB.Where("id = ?", query).First(&admins).Error
	if err == gorm.ErrRecordNotFound {
		global.Log.Warn("数据库菜单表中未发现查询的id")
		return
	} else if err != nil {
		global.Log.Warnf("数据库菜单表查询发生错误:%s", err)
		response.ResErrSrv(c, err)
		return
	}
	response.ResSuccess(c, &system.Admins{
		UserName:     admins.UserName,
		RealName:     admins.RealName,
		Email:        admins.Email,
		Phone:        admins.Phone,
		Status:       admins.Status,
		Avatar:       admins.Avatar,
		Introduction: admins.Introduction,
		DeptId:       admins.DeptId,
		PostId:       admins.PostId,
	})
}

// 获取用户下的角色ID列表
func GetAdminsRoleIDList(c *gin.Context) {
	query, b := c.GetQuery("adminsid")
	if !b {
		return
	}
	roleList := []uint64{}
	err := global.DB.Debug().Model(&system.AdminsRole{}).Where("admins_id = ?", query).Pluck("role_id", &roleList).Error
	if err != nil {
		response.ResErrSrv(c, err)
		return
	}
	response.ResSuccess(c, &roleList)
}

// AdminsUpdate 更新用户信息
func AdminsUpdate(c *gin.Context) {
	u1 := system.Admins{}
	err := c.ShouldBindJSON(&u1)
	if err != nil {
		fmt.Println("参数绑定出现问题")
		response.ResErrSrv(c, err)
		return
	}
	var out system.Admins
	err = global.DB.Where("user_name = ?", u1.UserName).First(&out).Error
	if err == gorm.ErrRecordNotFound { // record not found
		response.ResErrSrv(c, err)
		return
	} else if err != nil {
		response.ResErrSrv(c, err)
		return
	}
	err = global.DB.Debug().Model(&out).Updates(system.Admins{
		RealName: u1.RealName,
		Email:    u1.Email,
		Phone:    u1.Phone,
		Status:   u1.Status,
		DeptId:   u1.DeptId,
		PostId:   u1.PostId,
	}).Error
	if err != nil {
		response.ResFail(c, "更新失败...")
		return
	}
	response.ResSuccessMsg(c)
}

// AdminsCreate 新增用户
func AdminsCreate(c *gin.Context) {
	var admin system.Admins
	if err := c.ShouldBindJSON(&admin); err != nil {
		response.ResFail(c, "新增用户参数有误...")
		return
	}
	password, err := utils.EncryptPassword(admin.Password)
	if err != nil {
		response.ResFail(c, "用户密码加密失败")
		return
	}
	admin.Password = password
	err = global.DB.Create(&admin).Error
	if err != nil {
		response.ResErrSrv(c, err)
		return
	}
	response.ResSuccess(c, gin.H{"id": admin.ID})
}

// AdminsDelete 删除用户
func AdminsDelete(c *gin.Context) {
	var ids AdminsDelParam
	if err := c.ShouldBindJSON(&ids); err != nil {
		response.ResFail(c, "删除用户参数有误...")
		return
	}
	var uList []uint64
	for _, id := range ids.Id {
		uList = append(uList, id)
	}

	u1 := system.Admins{}
	err := u1.Delete(uList)
	if err != nil {
		response.ResErrSrv(c, err)
		return
	}
	response.ResSuccessMsg(c)

}

// 分配用户角色权限
func AdminsSetrole(c *gin.Context) {
	adminsid, b := c.GetQuery("adminsid")
	if !b {
		return
	}
	var roleId RoleIdRequestParam
	if err := c.ShouldBindJSON(&roleId); err != nil {
		response.ResFail(c, "删除用户参数有误...")
		return
	}
	intadminsid, err := strconv.Atoi(adminsid)
	if err != nil {
		return
	}
	uintadminsid := uint64(intadminsid) // adminsid uint64
	var roleids []uint64

	for _, param := range roleId.RoleId {
		roleids = append(roleids, param)
	}

	ar := system.AdminsRole{}
	err = ar.SetRole(uintadminsid, roleids)
	if err != nil {
		response.ResErrSrv(c, err)
		return
	}
	go core.CasbinAddRoleForUser(uintadminsid)
	response.ResSuccessMsg(c)
}
