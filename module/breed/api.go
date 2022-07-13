package breed

import (
	"net/http"
	"strings"

	"github.com/crazyhl/yzyx-materials/internal/db"
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
		message := err.Error()
		if strings.HasPrefix(err.Error(), "Error 1062:") {
			message = "该品种代码已存在，如列表没有，尝试去回收站恢复"
		}

		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": message,
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

	var form editForm
	if err := ctx.ShouldBind(&form); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "err: " + err.Error(),
		})
		return
	}

	breedDto, err := edit(ctx, form)
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
	err := delete(ctx)

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

type updateNetValueForm struct {
	NetValue float64 `form:"net_value" json:"net_value"`
}

// UpdateNetValue 更新净值
func UpdateNetValue(ctx *gin.Context) {
	var form updateNetValueForm
	if err := ctx.ShouldBind(&form); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "err: " + err.Error(),
		})
		return
	}

	breedDto, err := updateNetValue(ctx, form.NetValue)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "更新净值成功",
		"data":    breedDto,
	})
}

// AllList 全部列表，返回简单数据，只有 id code name
func AllList(ctx *gin.Context) {
	breeds := []*Breed{}
	breedDtos := []*BreedDto{}
	query := db.DB.Select([]string{"id", "code", "name"}).Where("user_id = ?", ctx.MustGet("user").(user.User).ID).Order("id desc")
	query.Find(&breeds)
	for _, breed := range breeds {
		breedDtos = append(breedDtos, breed.ToDto())
	}
	ctx.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "获取全部列表成功",
		"data":    breedDtos,
	})
}
