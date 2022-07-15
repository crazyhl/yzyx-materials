package route

import (
	"github.com/crazyhl/yzyx-materials/middlewares"
	"github.com/crazyhl/yzyx-materials/module/account"
	"github.com/crazyhl/yzyx-materials/module/breed"
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
			acc.GET("/list", account.List)
			acc.POST("/add", account.Add)
			acc.Use(account.CheckAccountUid())
			{
				acc.GET("/:id", account.Detail)
				acc.POST("/:id/bindBreed", account.BindBreed)
				acc.POST("/:id/addBreedBuyItem", account.AddBreedBuytItem)
				acc.DELETE("/delete/:id", account.Delete)
				acc.PUT("/update/:id", account.Edit)
			}
		}

		b := authorized.Group("/breed")
		{
			b.POST("/add", breed.Add)
			b.GET("/list", breed.List)
			b.GET("/allList", breed.AllList)
			b.Use(breed.CheckBreedUid)
			{
				b.PUT("/:id", breed.Edit)
				b.PUT("/:id/netValue", breed.UpdateNetValue)
				b.DELETE("/:id", breed.Delete)
			}
		}
	}
}
