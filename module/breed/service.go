package breed

import (
	"github.com/crazyhl/yzyx-materials/internal/db"
	"github.com/crazyhl/yzyx-materials/module/user"
	"github.com/gin-gonic/gin"
)

// add AddBreed 添加购买品种
func add(form AddBreedForm) (*BreedDto, error) {
	breed := &Breed{
		Code:      form.Code,
		Name:      form.Name,
		AccountId: form.AccountId,
	}
	if err := db.DB.Create(breed).Error; err != nil {
		return nil, err
	}

	return breed.ToDto(), nil
}

func delete(ctx *gin.Context, id uint) error {
	breed := &Breed{}
	db.DB.First(breed, id)
	if breed.Account.User.ID != ctx.MustGet("user").(user.User).ID {
		return ErrBreedNotYourAccount
	}

	return db.DB.Delete(&Breed{}, "id = ?", id).Error
}
