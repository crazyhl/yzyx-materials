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

			if claims, ok := token.Claims.(*user.UserJwtClaims); ok && token.Valid {
				fmt.Printf("%v ", claims)
			} else {
				fmt.Println(err)
			}

			c.Next()
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
