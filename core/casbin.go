package core

import (
	"gitee.com/go-server/global"
	"gitee.com/go-server/models/system"
	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"gorm.io/gorm"
	"strconv"
)

const (
	PrefixUserId = "u"
	PrefixRoleId = "r"
)

var Enforcer *casbin.Enforcer

// 角色-URL导入
func InitCasbinEnforcer(mysqldsn string) (err error) {
	var enforcer *casbin.Enforcer
	text := `[request_definition]
	r = sub, obj, act
	
	[policy_definition]
	p = sub, obj, act
	
	[role_definition]
	g = _, _
	
	[policy_effect]
	e = some(where (p.eft == allow))
	
	[matchers]
	m = g(r.sub, p.sub) == true \
			&& keyMatch2(r.obj, p.obj) == true \
			&& regexMatch(r.act, p.act) == true \
			|| r.sub == "root"`
	m, _ := model.NewModelFromString(text)
	a, _ := gormadapter.NewAdapter("mysql", mysqldsn, true) // Your driver and data source.
	enforcer, err = casbin.NewEnforcer(m, a)
	if err != nil {
		return err
	}
	//
	var roles []system.Role
	err = global.DB.Where(&system.Role{}).Find(&roles).Error
	if err != nil {
		return
	}
	if len(roles) == 0 {
		Enforcer = enforcer
		return
	}
	for _, role := range roles {
		// 设置角色权限
		setRolePermission(enforcer, role.ID)
	}
	Enforcer = enforcer
	return

}

// 设置角色权限
func setRolePermission(enforcer *casbin.Enforcer, id uint64) {
	var rolemenus []system.RoleMenu
	res := global.DB.Debug().Where("role_id = ?", id).Find(&rolemenus)
	if res.RowsAffected == 0 {
		global.Log.Fatal("没有查询到角色菜单数据......")
		return
	}
	// 查询角色菜单表，有数据 根据角色id来查询
	for _, rolemen := range rolemenus {
		var menu system.Menu // 菜单表
		err := global.DB.Model(&system.Menu{}).Where("id = ?", rolemen.MenuID).First(&menu).Error
		if err == gorm.ErrRecordNotFound { // record not found
			global.Log.Warnln("数据库中没有查询到结果record not found")
			return
		} else if err != nil {
			global.Log.Warnln("err里面有其他错误")
			return
		}
		if menu.MenuType == 3 { // 菜单类型 1模块 2菜单 3操作
			// 设置角色权限 Enforcer.AddPermissionsForUser("角色1","权限1", "访问规则")
			// 这里的访问规则 定为了GET/POST
			enforcer.AddPermissionForUser(PrefixRoleId+strconv.FormatUint(id, 10), menu.URL, "GET|POST")
		}
	}
}

// 设置用户角色
func CasbinSetRolePermission(roleid uint64) {
	if Enforcer == nil {
		return
	}
	// 删除用户角色
	Enforcer.DeletePermissionsForUser(PrefixRoleId + strconv.FormatUint(roleid, 10))
	// 设置用户角色
	setRolePermission(Enforcer, roleid)
}

// 删除角色
func CasbinDeleteRole(roleids []uint64) {
	if Enforcer == nil {
		return
	}
	for _, roleid := range roleids {
		// DeletePermissionsForUser 删除用户或角色的权限 Enforcer.DeletePermissionsForUser("bob")
		Enforcer.DeletePermissionsForUser(PrefixRoleId + strconv.FormatUint(roleid, 10))
		// DeleteRole 删除一个角色 Enforcer.DeleteRole("data2_admin")
		Enforcer.DeleteRole(PrefixRoleId + strconv.FormatUint(roleid, 10))
	}
}

// check user 是否有权限
func CasbinCheckPermission(userId, url, methodtype string) (bool, error) {
	return Enforcer.Enforce(PrefixUserId+userId, url, methodtype)
}

// 添加用户角色
func CasbinAddRoleForUser(userId uint64) (err error) {
	if Enforcer == nil {
		return
	}
	uid := PrefixUserId + strconv.FormatUint(userId, 10)
	Enforcer.DeleteRolesForUser(uid)
	var adminsroles []system.AdminsRole
	res := global.DB.Where("admins_id = ?", userId).Find(&adminsroles)
	if res.RowsAffected == 0 {
		global.Log.Warnln("用户角色表中没有找到结果")
		return
	}
	for _, adminsrole := range adminsroles {
		// 为用户添加角色 Enforcer.AddRoleForUser("用户1", "角色1")
		Enforcer.AddRoleForUser(uid, PrefixRoleId+strconv.FormatUint(adminsrole.RoleID, 10))
	}
	return
}
