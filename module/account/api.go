package account

import (
	"net/http"

	"github.com/crazyhl/yzyx-materials/internal/db"
	"github.com/crazyhl/yzyx-materials/module/domain/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm/clause"
)

type accountAddForm struct {
	Name               string  `form:"name" json:"name" binding:"required" label:"账户名称"`
	Description        string  `form:"description" json:"description" label:"账户描述"`
	ExpectTotalMoney   float64 `form:"expect_total_money" json:"expect_total_money" binding:"required_with:PerPartMoney|float|gt:0" label:"预计投入总金额"`
	PerPartMoney       float64 `form:"per_part_money" json:"per_part_money" binding:"float|gt:0" label:"每份投入金额"`
	ExpectRateOfReturn uint8   `form:"expect_rate_of_return" json:"expect_rate_of_return" binding:"uint|lte:100" label:"预计收益率"`
	User               models.User
}

// 添加账户
func Add(c *gin.Context) {
	var accAddForm accountAddForm
	if err := c.ShouldBind(&accAddForm); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "err: " + err.Error(),
		})
		return
	}
	accAddForm.User = c.MustGet("user").(models.User)
	accountDto, err := add(accAddForm)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "添加账户成功",
		"data":    accountDto,
	})
}

func List(c *gin.Context) {
	accounts := list(c)
	count := getCount(c)
	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "获取账户列表成功",
		"data": gin.H{
			"data":  accounts,
			"count": count,
		},
	})
}

func Delete(ctx *gin.Context) {
	err := delete(ctx)

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

type accountEditForm struct {
	Name               string  `form:"name" json:"name" binding:"required" label:"账户名称"`
	Description        string  `form:"description" json:"description" label:"账户描述"`
	ExpectTotalMoney   float64 `form:"expect_total_money" json:"expect_total_money" binding:"required_with:PerPartMoney|float|gt:0" label:"预计投入总金额"`
	PerPartMoney       float64 `form:"per_part_money" json:"per_part_money" binding:"float|gt:0" label:"每份投入金额"`
	ExpectRateOfReturn uint8   `form:"expect_rate_of_return" json:"expect_rate_of_return" binding:"uint|lte:100" label:"预计收益率"`
}

func Edit(ctx *gin.Context) {
	var accEditForm accountEditForm
	if err := ctx.ShouldBind(&accEditForm); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "err: " + err.Error(),
		})
		return
	}
	accountDto, err := edit(ctx, accEditForm)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "修改账户成功",
		"data":    accountDto,
	})
}

func Detail(c *gin.Context) {
	account := c.MustGet("account").(*models.Account)
	breeds := make([]*models.AccountBreed, 0)
	db.DB.Preload(clause.Associations).Where("account_id = ?", account.ID).Order("updated_at desc").Find(&breeds)
	account.Breeds = breeds

	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "获取账户成功",
		"data":    account.ToDto(),
	})
}
