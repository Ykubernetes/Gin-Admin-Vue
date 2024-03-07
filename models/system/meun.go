package system

import (
	"gitee.com/go-server/global"
	"gitee.com/go-server/models/basemodel"
	"github.com/ahmetb/go-linq"
	"gorm.io/gorm"
	"time"
)

// 菜单
type Menu struct {
	basemodel.Model
	Status      uint8  `gorm:"column:status;type:tinyint(1);not null;" json:"status" form:"status"`             // 状态(1:启用 2:不启用)
	Memo        string `gorm:"column:memo;size:64;" json:"memo" form:"memo"`                                    // 备注
	ParentID    uint64 `gorm:"column:parent_id;not null;" json:"parent_id" form:"parent_id"`                    // 父级ID
	URL         string `gorm:"column:url;size:72;" json:"url" form:"url"`                                       // 菜单URL
	Name        string `gorm:"column:name;size:32;not null;" json:"name" form:"name"`                           // 菜单名称
	Sequence    int    `gorm:"column:sequence;not null;" json:"sequence" form:"sequence"`                       // 排序值
	MenuType    uint8  `gorm:"column:menu_type;type:tinyint(1);not null;" json:"menu_type" form:"menu_type"`    // 菜单类型 1模块2菜单3操作
	Code        string `gorm:"column:code;size:32;not null;unique_index:uk_menu_code;" json:"code" form:"code"` // 菜单代码
	Icon        string `gorm:"column:icon;size:32;" json:"icon" form:"icon"`                                    // icon
	OperateType string `gorm:"column:operate_type;size:32;not null;" json:"operate_type" form:"operate_type"`   // 操作类型 none/add/del/view/update
}

// 添加前
func (m *Menu) BeforeCreate(tx *gorm.DB) error {
	m.CreatedAt = time.Now()
	m.UpdatedAt = time.Now()
	return nil
}

// 更新前
func (m *Menu) BeforeUpdate(tx *gorm.DB) error {
	m.UpdatedAt = time.Now()
	return nil
}

// 获取菜单有权限的操作列表
func (Menu) GetMenuButton(adminsid uint64, menuCode string, btns *[]string) (err error) {
	sql := `select operate_type from menu
	      where id in (
					select menu_id from role_menu where 
					menu_id in (select id from menu where parent_id in (select id from menu where code=?))
					and role_id in (select role_id from admins_role where admins_id=?)
				)`
	err = global.DB.Raw(sql, menuCode, adminsid).Pluck("operate_type", btns).Error
	return
}

// 获取管理员权限下所有菜单
func (Menu) GetMenuByAdminsid(adminsid uint64, menus *[]Menu) (err error) {
	sql := `select * from menu
	      where id in (
					select menu_id from role_menu where 
				  role_id in (select role_id from admins_role where admins_id=?)
				)`
	err = global.DB.Raw(sql, adminsid).Find(menus).Error
	return
}

// 删除菜单及关联数据
func (Menu) Delete(menuids []uint64) error {
	tx := global.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if err := tx.Error; err != nil {
		tx.Rollback()
		return err
	}
	for _, menuid := range menuids {
		if err := deleteMenuRecurve(tx, menuid); err != nil {
			tx.Rollback()
			return err
		}
	}
	if err := tx.Where("menu_id in (?)", menuids).Delete(&RoleMenu{}).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Where("id in (?)", menuids).Delete(&Menu{}).Error; err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}

// 递归删除
func deleteMenuRecurve(db *gorm.DB, parentID uint64) error {
	where := &Menu{}
	where.ParentID = parentID
	var menus []Menu
	dbslect := db.Where(&where)                        // 查询传送进来的parentID
	if err := dbslect.Find(&menus).Error; err != nil { // 把结果放到menus []Menu
		return err
	}
	for _, menu := range menus {
		if err := db.Where("menu_id = ?", menu.ID).Delete(&RoleMenu{}).Error; err != nil {
			return err
		}
		if err := deleteMenuRecurve(db, menu.ID); err != nil {
			return err
		}
	}
	if err := dbslect.Delete(&Menu{}).Error; err != nil {
		return err
	}
	return nil
}

// 查询所有菜单
func GetAllMenu() (menus []Menu, err error) {
	global.DB.Order("parent_id asc, sequence asc").Find(&menus)
	return
}

// 获取超级管理员初始菜单
func GetSuperAdminMenu() (out []MenuModel) {
	menuTop := MenuModel{
		Path:      "/sys",
		Component: "Sys",
		Name:      "Sys",
		Meta:      MenuMeta{Title: "系统管理", NoCache: false},
		Children:  []MenuModel{},
	}
	menuModel := MenuModel{
		Path:      "/icon",
		Component: "Icon",
		Name:      "Icon",
		Meta:      MenuMeta{Title: "图标管理", NoCache: false},
		Children:  []MenuModel{},
	}
	menuTop.Children = append(menuTop.Children, menuModel)
	menuModel = MenuModel{
		Path:      "/menu",
		Component: "Menu",
		Name:      "Menu",
		Meta:      MenuMeta{Title: "菜单管理", NoCache: false},
		Children:  []MenuModel{},
	}
	menuTop.Children = append(menuTop.Children, menuModel)
	menuModel = MenuModel{
		Path:      "role",
		Component: "Role",
		Name:      "Role",
		Meta:      MenuMeta{Title: "角色管理", NoCache: false},
		Children:  []MenuModel{},
	}
	menuTop.Children = append(menuTop.Children, menuModel)
	menuModel = MenuModel{
		Path:      "/admins",
		Component: "Admins",
		Name:      "Admins",
		Meta:      MenuMeta{Title: "用户管理", NoCache: false},
		Children:  []MenuModel{},
	}
	menuTop.Children = append(menuTop.Children, menuModel)
	out = append(out, menuTop)
	return
}

func SetMenu(menus []Menu, parentId uint64) (out []MenuModel) {
	var menuArr []Menu
	// 查询菜单表指定parentId的数据，以Sequence排序值排序，并返回结果到menuArr中
	linq.From(menus).Where(func(i interface{}) bool {
		return i.(Menu).ParentID == parentId
	}).OrderBy(func(i interface{}) interface{} {
		return i.(Menu).Sequence
	}).ToSlice(&menuArr)
	if len(menuArr) == 0 {
		return
	}
	noCache := false
	for _, item := range menuArr {
		menu := MenuModel{
			Path:      item.URL,
			Component: item.Code,
			Name:      item.Code,
			Meta:      MenuMeta{Title: item.Name, Icon: item.Icon, NoCache: noCache},
			Children:  []MenuModel{},
		}
		if item.MenuType == 3 {
			menu.Hidden = true
		}
		// 查询是否拥有子级
		menuChildren := SetMenu(menus, item.ID)
		if len(menuChildren) > 0 {
			menu.Children = menuChildren
		}
		if item.MenuType == 2 {
			// 添加子级别首页
			menuIndex := MenuModel{
				Path:      "index",
				Component: item.Code,
				Name:      item.Code,
				Meta:      MenuMeta{Title: item.Name, Icon: item.Icon, NoCache: noCache},
				Children:  []MenuModel{},
			}
			menu.Children = append(menu.Children, menuIndex)
			menu.Name = menu.Name + "index"
			menu.Meta = MenuMeta{}
		}
		out = append(out, menu)
	}
	return
}

// 查询登录用户权限菜单
func GetMenusByAdminsId(adminsid uint64) (result []Menu, err error) {
	menu := Menu{}
	var menus []Menu
	err = menu.GetMenuByAdminsid(adminsid, &menus)
	if err != nil || len(menus) == 0 {
		return
	}
	allmenu, err := GetAllMenu() // 所有菜单
	if err != nil || len(allmenu) == 0 {
		return
	}
	menuMapAll := make(map[uint64]Menu)
	for _, item := range allmenu {
		menuMapAll[item.ID] = item //
	}
	menuMap := make(map[uint64]Menu)
	for _, item := range menus {
		menuMap[item.ID] = item
	}
	for _, item := range menus {
		_, exists := menuMap[item.ParentID]
		if exists {
			continue
		}
		SetMenuUp(menuMapAll, item.ParentID, menuMap)
	}
	for _, m := range menuMap {
		result = append(result, m)
	}
	linq.From(result).OrderBy(func(i interface{}) interface{} {
		return i.(Menu).ParentID
	}).ToSlice(&result)
	return
}

// 向上查找父级菜单
func SetMenuUp(menuMapAll map[uint64]Menu, menuid uint64, menuMap map[uint64]Menu) {
	menuModel, exists := menuMapAll[menuid]
	if exists {
		mid := menuModel.ID
		_, exists = menuMap[mid]
		if !exists {
			menuMap[mid] = menuModel
			SetMenuUp(menuMapAll, menuModel.ParentID, menuMap)
		}
	}
}

// 新增菜单后自动添加菜单下的常规操作
func InitMenu(model Menu) {
	if model.MenuType != 2 {
		return
	}
	add := Menu{Status: 1, ParentID: model.ID, URL: model.URL + "/create", Name: "新增", Sequence: 1, MenuType: 3, Code: model.Code + "Add", OperateType: "add"}
	global.DB.Create(&add)
	del := Menu{Status: 1, ParentID: model.ID, URL: model.URL + "/delete", Name: "删除", Sequence: 2, MenuType: 3, Code: model.Code + "Del", OperateType: "del"}
	global.DB.Create(&del)
	view := Menu{Status: 1, ParentID: model.ID, URL: model.URL + "/detail", Name: "查看", Sequence: 3, MenuType: 3, Code: model.Code + "View", OperateType: "view"}
	global.DB.Create(&view)
	update := Menu{Status: 1, ParentID: model.ID, URL: model.URL + "/update", Name: "编辑", Sequence: 4, MenuType: 3, Code: model.Code + "Update", OperateType: "update"}
	global.DB.Create(&update)
	list := Menu{Status: 1, ParentID: model.ID, URL: model.URL + "/list", Name: "分页api", Sequence: 5, MenuType: 3, Code: model.Code + "List", OperateType: "list"}
	global.DB.Create(&list)
}
