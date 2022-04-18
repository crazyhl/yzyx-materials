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

func Login(ctx *gin.Context) {
	username := ctx.PostForm("username")
	password := ctx.PostForm("password")
	user, err := login(username, password)

	if err != nil {
		ctx.JSON(500, gin.H{
			"code":    500,
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(200, gin.H{
		"code":    200,
		"message": "登录成功",
		"data":    user,
	})
}
