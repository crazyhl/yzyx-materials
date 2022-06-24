package breed

import (
	"net/http"

	"github.com/crazyhl/yzyx-materials/internal/params"
	"github.com/crazyhl/yzyx-materials/module/user"
	"github.com/gin-gonic/gin"
)

type addForm struct {
	Code     string  `form:"code" json:"code" binding:"required"`
	Name     string  `form:"name" json:"name" binding:"required"`
	NetValue float64 `form:"net_value" json:"net_value"`
	Cost     float64 `form:"cost" json:"cost"`
	User     user.User
}

func Add(ctx *gin.Context) {
	var form addForm
	if err := ctx.ShouldBind(&form); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "err: " + err.Error(),
		})
		return
	}
	form.User = ctx.MustGet("user").(user.User)
	breedDto, err := add(form)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "添加品种成功",
		"data":    breedDto,
	})
}

type editForm struct {
	Code     string  `form:"code" json:"code" binding:"required"`
	Name     string  `form:"name" json:"name" binding:"required"`
	NetValue float64 `form:"net_value" json:"net_value"`
	Cost     float64 `form:"cost" json:"cost"`
}

func Edit(ctx *gin.Context) {
	id, err := params.GetUInt(ctx, "id")
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": err.Error(),
		})
		return
	}

	var form editForm
	if err := ctx.ShouldBind(&form); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "err: " + err.Error(),
		})
		return
	}

	breedDto, err := edit(form, id, ctx.MustGet("user").(user.User).ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "编辑品种成功",
		"data":    breedDto,
	})
}

func Delete(ctx *gin.Context) {
	id, err := params.GetUInt(ctx, "id")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "err: " + err.Error(),
		})
		return
	}

	err = delete(ctx.MustGet("user").(user.User).ID, id)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(200, gin.H{
		"code":    200,
		"message": "删除成功",
	})
}

func List(c *gin.Context) {
	breeds := list(c)
	count := getCount(c)
	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "获取品种列表成功",
		"data": gin.H{
			"data":  breeds,
			"count": count,
		},
	})
}
