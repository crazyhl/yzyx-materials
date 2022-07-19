package account

import (
	"net/http"

	"github.com/crazyhl/yzyx-materials/internal/params"
	"github.com/crazyhl/yzyx-materials/module/domain/models"
	"github.com/gin-gonic/gin"
)

// CheckAccountUid 账户操作相关中间件，检测账户是否存在以及是否是登录用户的账户
func CheckAccountUid() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := params.GetUInt(c, "id")
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"code":    http.StatusBadRequest,
				"message": "err: " + err.Error(),
			})
			return
		}
		account, err := GetByIdWithUidInternal(id, c.MustGet("user").(models.User).ID)

		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"code":    http.StatusBadRequest,
				"message": err.Error(),
			})
		} else {
			c.Set("account", account)
			c.Next()
		}
	}
}
