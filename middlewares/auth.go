package middlewares

import (
	"fmt"

	"github.com/crazyhl/yzyx-materials/module/user"
	"github.com/gin-gonic/gin"
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
			c.Abort()
		} else {
			fmt.Println(cliams)
			c.Next()
		}

	}
}
