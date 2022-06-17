package breed

import (
	"net/http"

	"github.com/crazyhl/yzyx-materials/module/user"
	"github.com/gin-gonic/gin"
)

type addForm struct {
	Code     string  `form:"code" json:"code" binding:"required"`
	Name     string  `form:"name" json:"name" binding:"required"`
	NetValue float64 `form:"net_value" json:"net_value" binding:"required"`
	Cost     float64 `form:"cost" json:"cost" binding:"required"`
	User     user.User
}

func Add(ctx gin.Context) {
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
