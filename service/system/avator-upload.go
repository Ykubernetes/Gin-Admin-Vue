package system

import (
	"fmt"
	"gitee.com/go-server/global"
	"gitee.com/go-server/models/system"
	"gitee.com/go-server/service/response"
	"gitee.com/go-server/utils"
	"github.com/gin-gonic/gin"
	"path"
	"strings"
)

func AvatorImgUpload(c *gin.Context) {
	// 从token 获取用户
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
	uname := userinfo.Username
	if uname == "" {
		response.ResFail(c, "获取用户信息失败")
		return
	}
	// 从请求中读取文件
	file, err := c.FormFile("avator_img")
	if err != nil {
		response.ResErrSrv(c, err)
		return
	} else {
		// 将读取到的文件保存在本地（服务端本地)
		dstFilePath := path.Join("upload/avator_img/", fmt.Sprintf("%s-%s", utils.GetUUID(), file.Filename))
		err = c.SaveUploadedFile(file, dstFilePath)
		if err != nil {
			response.ResFail(c, "图片上传失败")
			return
		} else {
			res := strings.Replace(dstFilePath, "\\", "/", -1)
			url_path := "http://" + c.Request.Host + "/" + res // 拼接全路径url
			resData := make(map[string]string)
			resData["img"] = url_path

			var u1 system.Admins
			global.DB.Where("user_name = ?", uname).First(&u1)
			err = global.DB.Model(&u1).Update("avatar", url_path).Error
			if err != nil {
				response.ResFail(c, "头像图片地址更新失败")
				return
			} else {
				response.ResSuccess(c, resData)
			}
		}
	}
}
