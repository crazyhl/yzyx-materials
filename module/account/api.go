package account

import (
	"net/http"

	"github.com/crazyhl/yzyx-materials/module/user"
	"github.com/gin-gonic/gin"
)

type accountAddForm struct {
	Name               string  `form:"name" json:"name" binding:"required" label:"账户名称"`
	Description        string  `form:"description" json:"description" label:"账户描述"`
	ExpectTotalMoney   float64 `form:"expect_total_money" json:"expect_total_money" binding:"float|gt:0" label:"预计投入总金额"`
	PerPartMoney       float64 `form:"per_part_money" json:"per_part_money" binding:"required_with:ExpectTotalMoney|float|gt:0" label:"每份投入金额"`
	ExpectRateOfReturn uint8   `form:"expect_rate_of_return" json:"expect_rate_of_return" binding:"uint|lte:100" label:"预计收益率"`
	User               user.User
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
	accAddForm.User = c.MustGet("user").(user.User)
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
	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "获取账户列表成功",
		"data":    accounts,
	})
}
