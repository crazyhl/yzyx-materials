package user

import (
	"github.com/gin-gonic/gin"
)

func Register(ctx *gin.Context) {
	username := ctx.PostForm("username")
	password := ctx.PostForm("password")
	_, err := register(username, password)

	if err != nil {
		ctx.JSON(500, gin.H{
			"code":    500,
			"message": "注册失败",
		})
		return
	}

	ctx.JSON(200, gin.H{
		"code":    200,
		"message": "注册成功",
	})
}
