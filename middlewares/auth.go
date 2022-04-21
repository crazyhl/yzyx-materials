package middlewares

import (
	"fmt"

	"github.com/crazyhl/yzyx-materials/module/user"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/spf13/viper"
)

func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		authorization := c.GetHeader("authorization")
		jwtStringHeader := authorization[0:6]
		if jwtStringHeader == "Bearer" {
			jwtStringBody := authorization[7:]
			fmt.Println(jwtStringBody)
			token, err := jwt.ParseWithClaims(jwtStringBody, &user.UserJwtClaims{}, func(token *jwt.Token) (interface{}, error) {
				return []byte(viper.GetString("JWT_SECRET")), nil
			})

			if err != nil {
				c.JSON(401, gin.H{
					"code":    401,
					"message": "校验失败，请重新登录",
				})
				c.Abort()
				return
			}

			if claims, ok := token.Claims.(*user.UserJwtClaims); ok && token.Valid {
				// 验证时间
				fmt.Println(claims.ExpiresAt)
				c.Next()
			} else {
				c.JSON(401, gin.H{
					"code":    401,
					"message": "校验失败，请重新登录2",
				})
				c.Abort()
			}

		} else {
			c.JSON(401, gin.H{
				"code":    401,
				"message": "请登录",
			})
			// 需要增加 abort 跳过后续执行
			c.Abort()
		}

	}
}
