package route

import (
	"github.com/crazyhl/yzyx-materials/middlewares"
	"github.com/crazyhl/yzyx-materials/module/account"
	"github.com/crazyhl/yzyx-materials/module/user"
	"github.com/gin-gonic/gin"
)

func InitRouter(router *gin.Engine) {
	router.POST("/user/register", user.Register)
	router.POST("/user/login", user.Login)

	authorized := router.Group("/")
	// AuthRequired() 中间件
	authorized.Use(middlewares.AuthRequired())
	{
		authorized.GET("/ping", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message": "pong",
			})
		})
		// account 接口
		acc := authorized.Group("/account")
		{
			acc.POST("/add", account.Add)
			acc.GET("/list", account.List)
			acc.DELETE("/delete/:id", account.Delete)
		}
	}
}
