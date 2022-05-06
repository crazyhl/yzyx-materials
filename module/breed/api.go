package breed

import (
	"net/http"
	"strconv"

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
		ctx.JSON(400, gin.H{
			"code":    400,
			"message": err.Error(),
		})
		return
	}
	// 根据账户ID查询账户信息
	account, err := account.GetByIdWithUidInternal(form.AccountId, ctx.MustGet("user").(user.User).ID)
	if err != nil {
		ctx.JSON(400, gin.H{
			"code":    400,
			"message": err.Error(),
		})
	}

	form.Account = *account
	breedDto, err := add(form)
	if err != nil {
		ctx.JSON(400, gin.H{
			"code":    400,
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

func Delete(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "err: " + err.Error(),
		})
		return
	}

	uintId := uint(id)
	err = delete(ctx, uintId)
	if err != nil {
		ctx.JSON(400, gin.H{
			"code":    400,
			"message": err.Error(),
		})
		return
	}
	ctx.JSON(200, gin.H{
		"code":    200,
		"message": "删除成功",
	})
}

type AddBreedItemForm struct {
	BreedID    uint    `form:"breed_id" binding:"required|uint" label:"品种ID"`
	Cost       float64 `form:"cost" binding:"required|float" label:"购买价格"`
	TotalPart  uint    `form:"total_part" binding:"required|uint" label:"总份数"`
	Commission float64 `form:"commission" binding:"required|float" label:"佣金"`
	Type       uint8   `form:"type" binding:"required|uint" label:"类型"`
}

func AddBuyItem(ctx *gin.Context) {
	var form AddBreedItemForm
	if err := ctx.ShouldBind(&form); err != nil {
		ctx.JSON(400, gin.H{
			"code":    400,
			"message": err.Error(),
		})
		return
	}
	// 验证 breed 所属账户所有权
	_, err := getBreedByIdWithUidInternal(form.BreedID, ctx.MustGet("user").(user.User).ID)
	if err != nil {
		ctx.JSON(400, gin.H{
			"code":    400,
			"message": err.Error(),
		})
		return
	}
	// 添加品种够买记录
	err = addBreedBuyItem(ctx, form)
	if err != nil {
		ctx.JSON(400, gin.H{
			"code":    400,
			"message": err.Error(),
		})
		return
	}
	// 更改品种各个值的计算

	ctx.JSON(200, gin.H{
		"code":    200,
		"message": "添加成功",
	})
}

func DeleteBuyItem(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "err: " + err.Error(),
		})
		return
	}

	uintId := uint(id)
	err = deleteBuyItem(ctx, uintId)
	if err != nil {
		ctx.JSON(400, gin.H{
			"code":    400,
			"message": err.Error(),
		})
		return
	}
	ctx.JSON(200, gin.H{
		"code":    200,
		"message": "删除成功",
	})
}

type EditBreedItemForm struct {
	Cost       float64 `form:"cost" binding:"required|float" label:"购买价格"`
	TotalPart  uint    `form:"total_part" binding:"required|uint" label:"总份数"`
	Commission float64 `form:"commission" binding:"required|float" label:"佣金"`
	Type       uint8   `form:"type" binding:"required|uint" label:"类型"`
}

func EditBuyItem(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "err: " + err.Error(),
		})
		return
	}
	uintId := uint(id)

	var form EditBreedItemForm
	if err := ctx.ShouldBind(&form); err != nil {
		ctx.JSON(400, gin.H{
			"code":    400,
			"message": err.Error(),
		})
		return
	}
	// 添加品种够买记录
	err = editBreedBuyItem(ctx, uintId, form)
	if err != nil {
		ctx.JSON(400, gin.H{
			"code":    400,
			"message": err.Error(),
		})
		return
	}
	ctx.JSON(200, gin.H{
		"code":    200,
		"message": "修改成功",
	})
}
