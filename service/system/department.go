package system

import (
	"gitee.com/go-server/global"
	"gitee.com/go-server/models/system"
	"gitee.com/go-server/service/response"
	"gitee.com/go-server/utils"
	"github.com/gin-gonic/gin"
	"strconv"
)

// GetDepartMentList 分页-部门列表数据
func GetDepartMentList(c *gin.Context) {
	var (
		Dept   system.Dept
		err    error
		result []system.Dept
	)
	Dept.DeptName = c.Request.FormValue("deptName")
	Dept.Status = c.Request.FormValue("status")
	Dept.ID, _ = strconv.ParseUint(c.Request.FormValue("deptId"), 10, 64)
	if Dept.DeptName == "" {
		result, err = Dept.SetDept(true)
	} else {
		result, err = Dept.GetPage(true)
	}
	if err != nil {
		response.ResFail(c, "获取部门列表失败")
		return
	}
	response.ResSuccess(c, result)
}

// GetDept 获取部门信息
func GetDept(c *gin.Context) {
	var (
		err  error
		Dept system.Dept
	)
	Dept.ID, _ = strconv.ParseUint(c.Param("id"), 10, 64)
	result, err := Dept.Get()
	if err != nil {
		response.ResFail(c, "获取部门信息出错..")
		return
	}
	response.ResSuccess(c, result)
}

// CreateDepartMent 添加部门
func CreateDepartMent(c *gin.Context) {
	var data system.Dept
	err := c.ShouldBindJSON(&data)
	if err != nil {
		response.ResFail(c, "添加部门,参数绑定失败...")
		return
	}
	// 创建人
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
		data.CreatedBy = 1
	} else {
		data.CreatedBy = userId
	}
	_, err = data.Create()
	if err != nil {
		response.ResFail(c, "创建部门失败...")
		return
	}
	response.ResSuccessMsg(c)
}

// UpdateDepartMent 修改部门
func UpdateDepartMent(c *gin.Context) {
	var data system.Dept
	err := c.ShouldBindJSON(&data)
	if err != nil {
		response.ResFail(c, "修改部门,参数绑定失败...")
		return
	}
	// 修改人
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
		data.UpdatedBy = 1
	} else {
		data.UpdatedBy = userId
	}

	_, err = data.Update(data.ID)
	if err != nil {
		response.ResFail(c, "修改部门失败...")
		return
	}
	response.ResSuccessMsg(c)
}

// DeleteDepartMent 删除部门
func DeleteDepartMent(c *gin.Context) {
	var data system.Dept
	Id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		global.Log.Warn("删除部门id获取失败,字符串转换int失败", err)
		return
	}
	_, err = data.Delete(Id)
	if err != nil {
		response.ResErrMsg(c, err)
		return
	}
	response.ResSuccessMsg(c)
}

// GetDeptTree 查询部门下拉树结构
func GetDeptTree(c *gin.Context) {
	var (
		Dept system.Dept
		err  error
	)
	Dept.DeptName = c.Request.FormValue("deptName")
	Dept.Status = c.Request.FormValue("status")

	deptId, _ := strconv.ParseUint(c.Request.FormValue("deptId"), 10, 64)

	Dept.ID = deptId
	result, err := Dept.SetDept(false)
	if err != nil {
		response.ResErrSrv(c, err)
		return
	}
	response.ResSuccess(c, result)
}
