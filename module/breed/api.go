package breed

import (
	"github.com/crazyhl/yzyx-materials/module/account"
	"github.com/crazyhl/yzyx-materials/module/user"
	"github.com/gin-gonic/gin"
)

// AddBreed 添加品种表单
type AddBreedForm struct {
	Code      string `form:"code" binding:"required" label:"品种编码"`
	Name      string `form:"name" binding:"required" label:"品种名称"`
	AccountId uint   `form:"account_id" binding:"required|uint" label:"账户ID"`
	Account   account.Account
}

func Add(ctx *gin.Context) {
	var form AddBreedForm
	if err := ctx.ShouldBind(&form); err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}
	// 根据账户ID查询账户信息
	account, err := account.GetByIdWithUidInternal(form.AccountId, ctx.MustGet("user").(user.User).ID)
	if err != nil {
		ctx.JSON(500, gin.H{
			"code":    500,
			"message": err.Error(),
		})
	}

	form.Account = *account
	breedDto, err := add(form)
	if err != nil {
		ctx.JSON(500, gin.H{
			"code":    500,
			"message": err.Error(),
		})
		return
	}
	ctx.JSON(200, gin.H{
		"code":    200,
		"message": "添加成功",
		"data":    breedDto,
	})
}
