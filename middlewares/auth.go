package middlewares

import (
	"github.com/crazyhl/yzyx-materials/module/user"
	"github.com/gin-gonic/gin"
	"github.com/golang-module/carbon/v2"
	log "github.com/sirupsen/logrus"
)

func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		authorization := c.GetHeader("authorization")
		cliams, err := user.ParseJwt(authorization)
		if err != nil {
			c.AbortWithStatusJSON(401, gin.H{
				"code":    401,
				"message": err.Error(),
			})
		} else {
			// 在这校验时间是否过期，用户是否存在
			now := carbon.Now().Timestamp()
			if cliams.ExpiresAt < now {
				c.AbortWithStatusJSON(401, gin.H{
					"code":    401,
					"message": user.ErrJWTExpired.Error(),
				})
				return
			}
			uid := cliams.ID
			loginUser, err := user.GetByUid(uid)
			if err != nil {
				log.Warn("user not found: ", err)
				c.AbortWithStatusJSON(401, gin.H{
					"code":    401,
					"message": user.ErrUserNotFound.Error(),
				})
				return
			}
			// 设置用户名为空，防止对外暴露
			loginUser.Password = ""
			// 把用户信息放到上下文中
			c.Set("user", loginUser)
			// 校验有效期，考虑是否续期，续期数据放到 header 中。
			compareExpiredAndRenewToken(c, now, cliams.ExpiresAt)
			c.Next()
		}
	}
}

// compareExpiredAndRenewToken 比较过期时间，如果过期，则续期
func compareExpiredAndRenewToken(c *gin.Context, now, claimsExpiresAt int64) {
	if now-claimsExpiresAt < 86400 {
		// 如果当前时间减去 jwt 的过期时间小于一天，则续期
		tokenStr, err := user.GenerateJWT(c.MustGet("user").(user.User))
		if err == nil {
			// 如果生成 token 成功，则把 token 放到 header 中
			c.Header("authorization", tokenStr)
		}
	}
}
