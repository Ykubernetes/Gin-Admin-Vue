package system

import (
	"gitee.com/go-server/global"
	"gitee.com/go-server/models/system"
	"gitee.com/go-server/service/response"
	"gitee.com/go-server/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"strconv"
)

type MenuListParam struct {
	Page      int    `form:"page" json:"page"`
	Limit     int    `form:"limit"`
	Sort      string `form:"sort"`
	Key       string `form:"key"`
	MenuType  uint   `form:"type"`
	Parent_id uint64 `form:"parent_id"`
}

type MenuIdDelParam struct {
	MenuId []uint64 `json:"menuid" binding:"required"`
}

// 获取所有菜单
func GetAllMenu(c *gin.Context) {
	var menus []system.Menu
	res := global.DB.Order("parent_id asc, sequence asc").Find(&menus)
	if res.RowsAffected == 0 {
		global.Log.Info("菜单表中查询结果为空")
		response.ResErrSrv(c, res.Error)
	}
	response.ResSuccess(c, &menus)
}

func GetMenuList(c *gin.Context) {
	var menuListParam MenuListParam
	err := c.ShouldBindQuery(&menuListParam)
	if err != nil {
		response.ResFail(c, "参数错误...")
		return
	}
	order := "ID DESC"
	var whereOrder []system.PageOrder

	if len(menuListParam.Sort) > 2 {
		orderType := menuListParam.Sort[0:1]
		order = menuListParam.Sort[1:len(menuListParam.Sort)]
		if orderType == "+" {
			order += " ASC"
		} else {
			order += " DESC"
		}
	}
	whereOrder = append(whereOrder, system.PageOrder{Order: order}) // [{id DESC  []}]

	if menuListParam.Key != "" {
		v := "%" + menuListParam.Key + "%"
		var arr []interface{}
		arr = append(arr, v)
		arr = append(arr, v)
		whereOrder = append(whereOrder, system.PageOrder{Where: "name like ? or code like ?", Value: arr})
	}
	if menuListParam.MenuType > 0 {
		var arr []interface{}
		arr = append(arr, menuListParam.MenuType)
		whereOrder = append(whereOrder, system.PageOrder{Where: "menu_type = ?", Value: arr})
	}
	if menuListParam.Parent_id > 0 {
		var arr []interface{}
		arr = append(arr, menuListParam.Parent_id)
		whereOrder = append(whereOrder, system.PageOrder{Where: "menu_type = ?", Value: arr})
	}
	var total int64
	list := []system.Menu{}
	table := global.DB.Table("menu")
	err = system.GetPage(table, &system.Menu{}, &list, menuListParam.Page, menuListParam.Limit, &total, whereOrder...)
	if err != nil {
		response.ResErrSrv(c, err)
		return
	}
	response.ResSuccessPage(c, total, &list)
}

// 详情
func Detail(c *gin.Context) {
	query, b := c.GetQuery("id")
	if !b {
		return
	}
	var menu system.Menu
	err := global.DB.Where("id = ?", query).First(&menu).Error
	if err == gorm.ErrRecordNotFound {
		global.Log.Warn("数据库菜单表中未发现查询的id")
		return
	} else if err != nil {
		global.Log.Warnf("数据库菜单表查询发生错误:%s", err)
		response.ResErrSrv(c, err)
		return
	}
	response.ResSuccess(c, &menu)
}

// 获取菜单有权限的操作列表
func MenuButtonList(c *gin.Context) {
	token := c.GetHeader("Authorization")
	if token == "" {
		response.ResFail(c, "token不存在！")
		return
	}
	userinfo, err := utils.ParseJwt(token, global.Config.JwtSecret.SecretKey)
	if err != nil {
		response.ResFail(c, "token解析失败")
		return
	}
	uid := global.RedisDB.Get(c, userinfo.UUID)
	userId, err := strconv.ParseUint(uid.Val(), 10, 64)
	if err != nil {
		global.Log.Warn("字符串转换uint64失败")
		return
	}
	menuCode, ok := c.GetQuery("menucode")
	if !ok {
		response.ResFail(c, "err")
		return
	}
	btnList := []string{}
	if userId == response.SUPER_ADMIN_ID {
		// 管理员
		btnList = append(btnList, "add")
		btnList = append(btnList, "del")
		btnList = append(btnList, "view")
		btnList = append(btnList, "update")
		btnList = append(btnList, "setrolemenu")
		btnList = append(btnList, "setadminrole")
	} else {
		menu := system.Menu{}
		err := menu.GetMenuButton(userId, menuCode, &btnList)
		if err != nil {
			response.ResErrSrv(c, err)
			return
		}
	}
	response.ResSuccess(c, &btnList)
}

// Create 菜单创建
func Create(c *gin.Context) {
	menu := system.Menu{}
	err := c.ShouldBindJSON(&menu)
	if err != nil {
		response.ResErrSrv(c, err)
		return
	}
	err = global.DB.Create(&menu).Error
	if err != nil {
		response.ResFail(c, "操作失败")
		return
	}
	go system.InitMenu(menu)
	response.ResSuccess(c, gin.H{"id": menu.ID})

}

// Delete 删除菜单
func Delete(c *gin.Context) {
	var list []uint64
	var menuids MenuIdDelParam
	if err := c.ShouldBindJSON(&menuids); err != nil {
		response.ResFail(c, "删除菜单参数有误...")
		return
	}
	for _, param := range menuids.MenuId {
		list = append(list, param)
	}

	menu := system.Menu{}
	err := menu.Delete(list) // 删除菜单及关联数据
	if err != nil {
		response.ResErrSrv(c, err)
		return
	}
	response.ResSuccessMsg(c)
}

// Update 删除
func Update(c *gin.Context) {
	menu := system.Menu{}
	err := c.ShouldBindJSON(&menu)
	if err != nil {
		response.ResErrSrv(c, err)
		return
	}
	err = global.DB.Save(&menu).Error
	if err != nil {
		response.ResFail(c, "操作失败")
		return
	}
	response.ResSuccessMsg(c)
}
