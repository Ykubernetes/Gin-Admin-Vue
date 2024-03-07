package middleware

import (
	"gitee.com/go-server/global"
	"gitee.com/go-server/service/response"
	"gitee.com/go-server/utils"
	"github.com/gin-gonic/gin"
	"time"
)

// 用户授权中间件
func UserAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		excludedRoutes := []string{
			"/user/login",
			"/user/logout",
			"/user/captcha/image",
			"/user/register",
			"/user/info",
			"/upload/*",
		}
		for _, excludedRoute := range excludedRoutes {
			if c.Request.URL.Path == excludedRoute {
				c.Next()
				return
			}
		}
		var uuid string
		token := c.GetHeader("Authorization")
		if token == "" {
			response.ResFail(c, "Token不存在，请重新登录")
			return
		}
		userinfo, err := utils.ParseJwt(token, global.Config.JwtSecret.SecretKey)
		if err != nil {
			response.ResFail(c, "token解析失败")
			return
		}
		exptimestamp := userinfo.ExpiresAt
		ok := exptimestamp.After(time.Now())
		if !ok {
			response.ResFailCode(c, "Token已过期,请重新登录", 5000)
			return
		} else {
			uuid = userinfo.UUID
		}

		if uuid == "" {
			response.ResFailCode(c, "用户未登录", 5000)
			return
		}
	}

}
