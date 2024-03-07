package system

import (
	"gitee.com/go-server/global"
	"gitee.com/go-server/service/response"
	"gitee.com/go-server/utils"
	"github.com/gin-gonic/gin"
)

func Getcaptcha(c *gin.Context) {
	code, s, err := utils.CreateCode()
	if err != nil {
		global.Log.Fatalf("图形验证码生成失败:%s", err)
		return
	}
	resData := make(map[string]string)
	resData["base64"] = s
	resData["key"] = code
	response.ResSuccess(c, &resData)
}
