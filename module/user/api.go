package user

import (
	"strings"

	"github.com/gin-gonic/gin"
)

func Register(ctx *gin.Context) {
	username := ctx.PostForm("username")
	password := ctx.PostForm("password")
	user, err := register(username, password)

	if err != nil {
		message := "注册失败"
		if strings.HasPrefix(err.Error(), "Error 1062:") {
			message += ": 用户名已存在"
		}

		ctx.JSON(500, gin.H{
			"code":    500,
			"message": message,
		})
		return
	}

	ctx.JSON(200, gin.H{
		"code":    200,
		"message": "注册成功",
		"data":    user,
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
