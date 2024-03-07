package middleware

import (
	"fmt"
	"gitee.com/go-server/core"
	"gitee.com/go-server/global"
	"gitee.com/go-server/service/response"
	"gitee.com/go-server/utils"
	"github.com/gin-gonic/gin"
	"strconv"
)

func CasbinCheckMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		excludedRoutes := []string{
			"/user/login",
			"/user/logout",
			"/user/captcha/image",
			"/menu/menubuttonlist",
			"/menu/allmenu",
			"/admins/adminsroleidlist",
			"/user/info",
			"/user/systemstate", // 公共的系统负载接口 跳过检查权限，不跳过检查auth
			"/user/register",
			"/user/editpwd",
			"/role/rolemenuidlist",
			"/role/allrole",
			"/upload/*",
			"/profile/index",
			"/user/avatar/img",
			"/user/editinfo",
			"/user/systeminfo",
		}
		for _, excludedRoute := range excludedRoutes {
			if c.Request.URL.Path == excludedRoute {
				c.Next()
				return
			}
		}
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
			c.Next()
			return
		}
		p := c.Request.URL.Path
		m := c.Request.Method
		useid := uid.Val()
		if b, err := core.CasbinCheckPermission(useid, p, m); err != nil {
			response.ResFail(c, "err303"+err.Error())
			fmt.Println("err303**", err)
			return
		} else if !b {
			response.ResFail(c, "您没有访问权限")
			return
		}
		c.Next()
	}

}
