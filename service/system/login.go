package system

import (
	"gitee.com/go-server/core"
	"gitee.com/go-server/global"
	"gitee.com/go-server/models/system"
	"gitee.com/go-server/service/response"
	"gitee.com/go-server/utils"
	"github.com/gin-gonic/gin"
	"strconv"
)

// LoginParams bind request params
type LoginParams struct {
	UserName string `json:"username" binding:"required"`
	PassWord string `json:"password" binding:"required"`
	PicKey   string `json:"picKey" binding:"required" `
	Code     string `json:"code" binding:"required"`
}

// LoginChangePassWordParams 修改密码请求参数
type LoginChangePassWordParams struct {
	OpassWord      string `json:"oldPassword" binding:"required"`
	NpassWord      string `json:"newPassword" binding:"required"`
	NpassWordAgain string `json:"newPassword_again" binding:"required"`
}

// LoginEditInfoParams 个人中心修改信息请求参数
type LoginEditInfoParams struct {
	Email        string `json:"email" binding:"required"`
	Introduction string `json:"introduction" binding:"required"`
	Realname     string `json:"realname" binding:"required"`
}

// 这个注册暂时不能使用
func Register(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")
	if username == "" || password == "" {
		response.ResFail(c, "用户名或密码不能为空")
		return
	}
	if len(username) < 4 || len(password) < 6 {
		response.ResFail(c, "用户名或密码太简单")
		return
	}
	encryptPassword, err := utils.EncryptPassword(password)
	if err != nil {
		global.Log.Warnln("注册密码加密失败", err)
		return
	}
	user := system.Admins{
		UserName:     username,
		RealName:     "",
		Password:     encryptPassword,
		Email:        "admin@shenfor.com.cn",
		Phone:        "13813808888",
		Status:       1,
		Avatar:       "https://www.telnote.cn/uploads/touxiang/80/014126751420.jpeg",
		Introduction: "这个家伙很懒，什么都没有留下",
	}
	global.DB.Create(&user)
	response.ResSuccessMsg(c)

}

func Login(c *gin.Context) {
	var loginParam LoginParams
	if err := c.ShouldBindJSON(&loginParam); err != nil {
		response.ResFail(c, "登录参数有误...")
		return
	}
	username := loginParam.UserName
	password := loginParam.PassWord
	captchaKey := loginParam.PicKey
	verifyValue := loginParam.Code
	if username == "" || password == "" {
		response.ResFail(c, "用户名或密码不能为空")
		return
	}
	userinfo, err := system.FindByUsername(username)
	if err != nil {
		response.ResFail(c, "用户不存在")
		return
	}
	val := utils.VerifyCaptcha(captchaKey, verifyValue)
	if val {
		if !utils.EqualsPassword(password, userinfo.Password) {
			response.ResFail(c, "用户名或密码错误")
			return
		}
		if username == userinfo.UserName && username == "admin" {
			//
			userinfo.ID = response.SUPER_ADMIN_ID
			userinfo.Status = 1 // 1:正常 2:未激活 3:暂停使用
		}
		if userinfo.Status != 1 {
			response.ResFail(c, "该用户已被禁用，请联系管理员")
			return
		}

		// redis
		uuid := utils.GetUUID()
		ok := global.RedisDB.Set(c, uuid, userinfo.ID, 0).Val()
		if ok != "OK" {
			global.Log.Fatal("Redis存储用户信息出错...")
		}
		// token
		token, err := utils.GenerateJWT(userinfo.UserName, uuid, global.Config.JwtSecret.SecretKey)
		if err != nil {
			global.Log.Fatal("Token生成失败...")
			return
		}
		// 返回token
		resData := make(map[string]string)
		resData["token"] = token

		// casbin 处理
		err = core.CasbinAddRoleForUser(userinfo.ID)
		if err != nil {
			response.ResErrSrv(c, err)
			return
		}
		response.ResSuccess(c, resData)
	} else {
		response.ResFail(c, "验证码错误...")
	}

}

func Logout(c *gin.Context) {
	token := c.GetHeader("Authorization")
	if token == "" {
		response.ResFail(c, "操作失败1")
		return
	}
	userinfo, err := utils.ParseJwt(token, global.Config.JwtSecret.SecretKey)
	if err != nil {
		response.ResFail(c, "操作失败2")
		return
	}
	cid := userinfo.UUID
	if cid == "" {
		response.ResFail(c, "操作失败3")
		return
	}
	// 清除redis中的userid 它的key 是uuid生成的 值是user.id
	global.RedisDB.Del(c, cid)
	response.ResSuccessMsg(c)
}

// 获取用户信息及可访问的权限菜单
func GetUserInfo(c *gin.Context) {
	var menuData []system.Menu
	var err error
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
	if userId == response.SUPER_ADMIN_ID {
		// 管理员
		menuData, err = system.GetAllMenu()
		if err != nil {
			response.ResErrSrv(c, err)
			return
		}
		if len(menuData) == 0 {
			menuModelTop := system.Menu{
				Status:      1,
				ParentID:    0,
				URL:         "",
				Name:        "TOP",
				Sequence:    1,
				MenuType:    1,
				Code:        "TOP",
				OperateType: "none",
			}
			global.DB.Create(&menuModelTop)
			menuModelSys := system.Menu{
				Status:      1,
				ParentID:    menuModelTop.ID,
				URL:         "",
				Name:        "系统管理",
				Sequence:    1,
				MenuType:    1,
				Code:        "Sys",
				Icon:        "el-icon-menu",
				OperateType: "none",
			}
			global.DB.Create(&menuModelSys)
			menuModel := system.Menu{
				Status:      1,
				ParentID:    menuModelSys.ID,
				URL:         "/menu",
				Name:        "菜单管理",
				Sequence:    20,
				MenuType:    2,
				Code:        "Menu",
				Icon:        "documentation",
				OperateType: "none",
			}
			global.DB.Create(&menuModel)
			system.InitMenu(menuModel)
			menuModel = system.Menu{Status: 1, ParentID: menuModelSys.ID, URL: "/role", Name: "角色管理", Sequence: 30, MenuType: 2, Code: "Role", Icon: "tree", OperateType: "none"}
			global.DB.Create(&menuModel)
			system.InitMenu(menuModel)
			menuModel = system.Menu{Status: 1, ParentID: menuModelSys.ID, URL: "/department", Name: "部门管理", Sequence: 30, MenuType: 2, Code: "Dept", Icon: "el-icon-s-order", OperateType: "none"}
			global.DB.Create(&menuModel)
			system.InitMenu(menuModel)
			menuModel = system.Menu{Status: 1, ParentID: menuModelSys.ID, URL: "/post", Name: "岗位管理", Sequence: 30, MenuType: 2, Code: "Post", Icon: "el-icon-office-building", OperateType: "none"}
			global.DB.Create(&menuModel)
			system.InitMenu(menuModel)
			menuModel = system.Menu{Status: 1, ParentID: menuModel.ID, URL: "/role/setrole", Name: "分配角色菜单", Sequence: 6, MenuType: 3, Code: "RoleSetrolemenu", Icon: "", OperateType: "setrolemenu"}
			global.DB.Create(&menuModel)
			menuModel = system.Menu{Status: 1, ParentID: menuModelSys.ID, URL: "/admins", Name: "用户管理", Sequence: 40, MenuType: 2, Code: "Admins", Icon: "user", OperateType: "none"}
			global.DB.Create(&menuModel)
			system.InitMenu(menuModel)
			menuModel = system.Menu{Status: 1, ParentID: menuModel.ID, URL: "/admins/setrole", Name: "分配角色", Sequence: 6, MenuType: 3, Code: "AdminsSetrole", Icon: "", OperateType: "setadminrole"}
			global.DB.Create(&menuModel)
			menuData, _ = system.GetAllMenu()
		}
	} else {
		menuData, err = system.GetMenusByAdminsId(userId)
		if err != nil {
			response.ResErrSrv(c, err)
			return
		}
	}

	var menus []system.MenuModel
	if len(menuData) > 0 {
		var topmenuid uint64 = menuData[0].ParentID
		if topmenuid == 0 {
			topmenuid = menuData[0].ID
		}
		menus = system.SetMenu(menuData, topmenuid)
	}
	if len(menus) == 0 && userId == response.SUPER_ADMIN_ID {
		menus = system.GetSuperAdminMenu()
	}
	userinfo2, err := system.FindByUserId(userId)
	if err != nil {
		response.ResFail(c, "获取用户信息失败...")
		return
	}
	resData := system.UserData{
		Menus:        menus,
		Introduction: userinfo2.Introduction,
		Avatar:       userinfo2.Avatar,
		Name:         userinfo2.UserName,
		Email:        userinfo2.Email,
	}
	response.ResSuccess(c, &resData)
}

func ChangePassword(c *gin.Context) {
	var loginPassWDParams LoginChangePassWordParams
	token := c.GetHeader("Authorization")
	if token == "" {
		response.ResFail(c, "操作失败1")
		return
	}
	userinfo, err := utils.ParseJwt(token, global.Config.JwtSecret.SecretKey)
	if err != nil {
		response.ResFail(c, "操作失败2")
		return
	}
	username := userinfo.Username
	if username == "" {
		response.ResFail(c, "操作失败3")
		return
	}

	if err := c.ShouldBindJSON(&loginPassWDParams); err != nil {
		response.ResFail(c, "密码修改参数有误...")
		return
	}

	new_password := loginPassWDParams.NpassWord
	new_password_again := loginPassWDParams.NpassWordAgain
	if len(new_password) < 6 || len(new_password) > 20 {
		response.ResFail(c, "密码长度在 6 到 20 个字符")
		return
	}
	if new_password != new_password_again {
		response.ResFail(c, "两次密码不匹配")
		return
	}
	var u1 system.Admins
	global.DB.Where("user_name = ?", username).First(&u1)
	old_password := loginPassWDParams.OpassWord
	if !utils.EqualsPassword(old_password, u1.Password) {
		response.ResFail(c, "原密码不匹配")
		return
	}
	encryptPassword, err := utils.EncryptPassword(new_password)
	if err != nil {
		global.Log.Warnln("密码加密失败", err)
		return
	}
	err = global.DB.Model(&u1).Update("password", encryptPassword).Error
	if err != nil {
		response.ResFail(c, "密码更新失败")
		return
	}
	response.ResSuccessMsg(c)
}

func ChangeUserinfo(c *gin.Context) {
	var loginEditParams LoginEditInfoParams
	token := c.GetHeader("Authorization")
	if token == "" {
		response.ResFail(c, "操作失败1")
		return
	}
	userinfo, err := utils.ParseJwt(token, global.Config.JwtSecret.SecretKey)
	if err != nil {
		response.ResFail(c, "操作失败2")
		return
	}
	username := userinfo.Username
	if username == "" {
		response.ResFail(c, "操作失败3")
		return
	}

	if err := c.ShouldBindJSON(&loginEditParams); err != nil {
		response.ResFail(c, "密码修改参数有误...")
		return
	}

	new_email := loginEditParams.Email
	new_introduction := loginEditParams.Introduction
	new_realname := loginEditParams.Realname
	var u1 system.Admins
	global.DB.Where("user_name = ?", username).First(&u1)
	if new_email == u1.Email && new_introduction == u1.Introduction && new_realname == u1.RealName {
		response.ResFail(c, "未做任何修改")
		return
	}
	if new_email == "" || new_introduction == "" {
		response.ResFail(c, "请修改信息后再提交")
		return
	}
	err = global.DB.Model(&u1).Updates(system.Admins{
		RealName:     new_realname,
		Email:        new_email,
		Introduction: new_introduction,
	}).Error
	if err != nil {
		response.ResFail(c, "用户信息更新失败")
		return
	}
	response.ResSuccessMsg(c)
}
