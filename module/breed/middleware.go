package breed

import (
	"net/http"

	"github.com/crazyhl/yzyx-materials/internal/params"
	"github.com/crazyhl/yzyx-materials/module/user"
	"github.com/gin-gonic/gin"
)

func CheckBreedUid(ctx *gin.Context) {
	id, err := params.GetUInt(ctx, "id")
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "err: " + err.Error(),
		})
		return
	}
	breed, err := GetByIdWithUidInternal(id, ctx.MustGet("user").(user.User).ID)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "err: " + err.Error(),
		})
	} else {
		ctx.Set("breed", breed)
		ctx.Next()
	}
}
