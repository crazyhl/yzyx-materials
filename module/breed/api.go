package breed

import (
	"github.com/crazyhl/yzyx-materials/module/account"
	"github.com/gin-gonic/gin"
)

type AddBreedForm struct {
	Code    string `form:"code" binding:"required"`
	Name    string `form:"name" binding:"required"`
	Account account.Account
}

func Add(ctx *gin.Context) {

}
