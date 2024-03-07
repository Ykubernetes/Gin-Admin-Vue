package system

import (
	"gitee.com/go-server/global"
	"gitee.com/go-server/models/system"
	"gitee.com/go-server/service/response"
	"gitee.com/go-server/utils"
	"github.com/gin-gonic/gin"
	"strconv"
)

// GetPostList 岗位列表
func GetPostList(c *gin.Context) {
	var (
		data      system.Post
		err       error
		pageSize  = 10
		pageIndex = 1
	)
	if size := c.Request.FormValue("pageSize"); size != "" {
		result, err := strconv.Atoi(size)
		if err != nil {
			global.Log.Warnln("pageSize string 转换int 发生错误", err)
			return
		}
		pageSize = result
	}
	if index := c.Request.FormValue("pageIndex"); index != "" {
		result, err := strconv.Atoi(index)
		if err != nil {
			global.Log.Warnln("pageIndex string 转换int 发生错误", err)
			return
		}
		pageIndex = result
	}
	postId := c.Request.FormValue("id")
	if postId != "" {
		atoi, err := strconv.ParseUint(postId, 10, 64)
		if err != nil {
			global.Log.Warnln("postId string 转换uint64 发生错误", err)
			return
		}
		data.ID = atoi
	}

	data.PostCode = c.Request.FormValue("postCode")
	data.PostName = c.Request.FormValue("postName")
	data.Status = c.Request.FormValue("status")

	result, count, err := data.GetPage(pageSize, pageIndex)
	if err != nil {
		global.Log.Fatalln("获取岗位列表失败,Err:", err)
		response.ResFail(c, "获取岗位列表失败...")
		return
	}
	response.ResSuccessPage(c, count, &result)
}

func GetPost(c *gin.Context) {
	var (
		err  error
		Post system.Post
	)
	atoi, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		global.Log.Warnln("postId string转换int,发生错误：", err)
		return
	}
	Post.ID = atoi

	result, err := Post.Get()
	if err != nil {
		response.ResFail(c, "获取岗位信息失败...")
		return
	}
	response.ResSuccess(c, result)
}

// CreatePost 添加岗位
func CreatePost(c *gin.Context) {
	var data system.Post
	err := c.ShouldBindJSON(&data)
	if err != nil {
		global.Log.Warnln("岗位参数绑定失败:", err)
		response.ResFail(c, "岗位参数绑定失败")
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

	result, err := data.Create()
	if err != nil {
		global.Log.Warnln("岗位添加失败:", err)
		response.ResFail(c, "岗位添加失败")
		return
	}
	response.ResSuccess(c, result)
}

// UpdatePost 修改岗位
func UpdatePost(c *gin.Context) {
	var data system.Post
	err := c.ShouldBindJSON(&data)
	if err != nil {
		global.Log.Warnln("绑定岗位参数失败:", err)
		response.ResFail(c, "绑定岗位参数失败")
		return
	}
	// 更新人
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
	result, err := data.Update(data.ID)
	if err != nil {
		global.Log.Fatalln("更新岗位信息失败", err)
		response.ResFail(c, "更新岗位信息失败")
		return
	}
	response.ResSuccess(c, result)
}

// DeletePost 删除岗位
func DeletePost(c *gin.Context) {
	var data system.Post
	// 更新人
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
	IDS := utils.IdsStrToIdsIntGroup("id", c)
	result, err := data.BatchDelete(IDS)
	if err != nil {
		global.Log.Fatalln("删除岗位失败:", err)
		response.ResFail(c, "删除岗位失败...")
		return
	}
	response.ResSuccess(c, result)
}
