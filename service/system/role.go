package system

import (
	"gitee.com/go-server/core"
	"gitee.com/go-server/global"
	"gitee.com/go-server/models/system"
	"gitee.com/go-server/service/response"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"strconv"
)

type RolesListParam struct {
	Page      int    `form:"page" json:"page"`
	Limit     int    `form:"limit"`
	Sort      string `form:"sort"`
	Key       string `form:"key"`
	Parent_id uint64 `form:"parent_id"`
}

// 分页条件
type PageWhereOrder struct {
	Order string
	Where string
	Value []interface{}
}

type RoleIdDelParam struct {
	RoleId []uint64 `json:"roleid" binding:"required"`
}

type RolePermissionParam struct {
	PPM []uint64 `json:"ppm" binding:"required"`
}

// 所有角色
func GetAllRole(c *gin.Context) {
	var list []system.Role
	err := global.DB.Find(&list).Order("parent_id asc,sequence asc").Error
	if err != nil {
		response.ResErrSrv(c, err)
		return
	}
	response.ResSuccess(c, &list)
}

// 分页数据
func GetRoleList(c *gin.Context) {
	var rolesListParam RolesListParam
	err := c.ShouldBindQuery(&rolesListParam)
	if err != nil {
		response.ResFail(c, "参数绑定错误...")
		return
	}
	var whereOrder []system.PageOrder
	order := "ID DESC"
	if len(rolesListParam.Sort) >= 2 {
		orderType := rolesListParam.Sort[0:1]
		order = rolesListParam.Sort[1:len(rolesListParam.Sort)]
		if orderType == "+" {
			order += " ASC"
		} else {
			order += " DESC"
		}
	}
	whereOrder = append(whereOrder, system.PageOrder{Order: order})
	if rolesListParam.Key != "" {
		v := "%" + rolesListParam.Key + "%"
		var arr []interface{}
		arr = append(arr, v)
		whereOrder = append(whereOrder, system.PageOrder{Where: "name like ?", Value: arr})
	}
	if rolesListParam.Parent_id > 0 {
		var arr []interface{}
		arr = append(arr, rolesListParam.Parent_id)
		whereOrder = append(whereOrder, system.PageOrder{Where: "parent_id = ?", Value: arr})
	}
	var total int64
	list := []system.Role{}
	table := global.DB.Table("role")
	err = system.GetPage(table, &system.Role{}, &list, rolesListParam.Page, rolesListParam.Limit, &total, whereOrder...)
	if err != nil {
		response.ResErrSrv(c, err)
		return
	}
	response.ResSuccessPage(c, total, &list)
}

// 获取角色详情
func GetRoleDetail(c *gin.Context) {
	query, b := c.GetQuery("id")
	if !b {
		return
	}
	var role system.Role
	err := global.DB.Where("id = ?", query).First(&role).Error
	if err == gorm.ErrRecordNotFound {
		global.Log.Warn("数据库菜单表中未发现查询的id")
		return
	} else if err != nil {
		global.Log.Warnf("数据库菜单表查询发生错误:%s", err)
		response.ResErrSrv(c, err)
		return
	}
	response.ResSuccess(c, &role)
}

// 获取角色下的菜单ID列表
func GetRoleMenuIDList(c *gin.Context) {
	query, b := c.GetQuery("roleid")
	if !b {
		return
	}
	roleList := []uint64{}
	err := global.DB.Debug().Model(&system.RoleMenu{}).Where("role_id = ?", query).Pluck("menu_id", &roleList).Error
	if err != nil {
		response.ResErrSrv(c, err)
		return
	}
	response.ResSuccess(c, &roleList)

}

// RoleSetRole 设置角色菜单权限
func RoleSetRole(c *gin.Context) {
	var ppm RolePermissionParam
	if err := c.ShouldBindJSON(&ppm); err != nil {
		response.ResFail(c, "设置角色权限参数有误...")
		return
	}
	query, b := c.GetQuery("roleid")
	if !b {
		return
	}
	atoi, err := strconv.Atoi(query)
	if err != nil {
		return
	}
	uquery := uint64(atoi)
	var list []uint64

	for _, param := range ppm.PPM {
		list = append(list, param)
	}
	rm := system.RoleMenu{}
	err = rm.SetRole(uquery, list)
	if err != nil {
		response.ResErrSrv(c, err)
		return
	}
	go core.CasbinSetRolePermission(uquery)
	response.ResSuccessMsg(c)
}

// RoleCreate 角色创建
func RoleCreate(c *gin.Context) {
	r1 := system.Role{}
	err := c.ShouldBindJSON(&r1)
	if err != nil {
		response.ResErrSrv(c, err)
		return
	}
	err = global.DB.Create(&r1).Error
	if err != nil {
		response.ResErrSrv(c, err)
		return
	}
	response.ResSuccess(c, gin.H{"id": r1.ID})
}

// RoleUpdate 编辑角色
func RoleUpdate(c *gin.Context) {
	r1 := system.Role{}
	err := c.ShouldBindJSON(&r1)
	if err != nil {
		response.ResErrSrv(c, err)
		return
	}
	err = global.DB.Save(&r1).Error
	if err != nil {
		response.ResFail(c, "角色更新失败...")
		return
	}
	response.ResSuccessMsg(c)
}

// RoleDelete 角色删除
func RoleDelete(c *gin.Context) {
	var list []uint64
	var roleids RoleIdDelParam
	if err := c.ShouldBindJSON(&roleids); err != nil {
		response.ResFail(c, "删除角色参数有误...")
		return
	}
	for _, param := range roleids.RoleId {
		list = append(list, param)
	}
	role := system.Role{}
	err := role.Delete(list)
	if err != nil {
		response.ResErrSrv(c, err)
		return
	}
	go core.CasbinDeleteRole(list)
	response.ResSuccessMsg(c)
}
