package accountbreed

import (
	"net/http"

	"github.com/crazyhl/yzyx-materials/internal/params"
	"github.com/crazyhl/yzyx-materials/module/breed"
	"github.com/crazyhl/yzyx-materials/module/domain/models"
	"github.com/gin-gonic/gin"
)

// BindBreed 绑定品种
func BindBreed(ctx *gin.Context) {
	// 此处通过校验中间件已经可以成功获取到 account 了
	// 然后就使用 form 接收 breed id 进行校验就可以了
	breedDto, err := bindBreed(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "err: " + err.Error(),
		})
	} else {
		ctx.JSON(http.StatusOK, gin.H{
			"code":    http.StatusOK,
			"message": "绑定品种成功",
			"data":    breedDto,
		})
	}
}

func AddBreedBuytItem(ctx *gin.Context) {
	breedDto, err := addBreedBuytItem(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "err: " + err.Error(),
		})
	} else {
		ctx.JSON(http.StatusOK, gin.H{
			"code":    http.StatusOK,
			"message": "添加购买记录成功",
			"data":    breedDto,
		})
	}
}

func UpdateBreedBuytItem(ctx *gin.Context) {
	breedItemDto, err := updateBreedBuytItem(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "err: " + err.Error(),
		})
	} else {
		ctx.JSON(http.StatusOK, gin.H{
			"code":    http.StatusOK,
			"message": "添加购买记录成功",
			"data":    breedItemDto,
		})
	}
}

func BuyItemList(ctx *gin.Context) {
	breedId, err := params.GetUInt(ctx, "id")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "err: " + err.Error(),
		})
	}
	b, err := breed.GetByIdWithUidInternal(breedId, ctx.MustGet("user").(models.User).ID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "err: " + err.Error(),
		})
	}
	acc := ctx.MustGet("account").(*models.Account)
	buyItems := buyItemList(ctx)
	count := getBuyItemCount(ctx)
	ctx.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "账户: " + acc.Name + " : " + b.Name + "(" + b.Code + ")购买记录",
		"data": gin.H{
			"data":  buyItems,
			"count": count,
		},
	})
}
