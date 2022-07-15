package account

import (
	"github.com/crazyhl/yzyx-materials/internal/db"
	"github.com/crazyhl/yzyx-materials/internal/model"
	"github.com/crazyhl/yzyx-materials/module/breed"
	"github.com/crazyhl/yzyx-materials/module/user"
	"github.com/gin-gonic/gin"
	"github.com/golang-module/carbon/v2"
)

// addAccount 添加账户
func add(form accountAddForm) (*AccountDto, error) {
	account := &Account{
		Name:               form.Name,
		Description:        form.Description,
		User:               form.User,
		ExpectTotalMoney:   form.ExpectTotalMoney,
		PerPartMoney:       form.PerPartMoney,
		ExpectRateOfReturn: form.ExpectRateOfReturn,
	}
	if err := db.DB.Create(account).Error; err != nil {
		return nil, err
	}

	// 将 account 转换为 AccountDto
	accountDto := account.ToDto()

	return accountDto, nil
}

// listAccounts 获取账户列表
func list(c *gin.Context) []*AccountDto {
	accounts := []*Account{}
	accountDtos := []*AccountDto{}
	db.DB.Scopes(db.Paginate(c)).Where("user_id = ?", c.MustGet("user").(user.User).ID).Order("id desc").Find(&accounts)
	for _, account := range accounts {
		accountDtos = append(accountDtos, account.ToDto())
	}
	return accountDtos
}

func getCount(c *gin.Context) int64 {
	count := int64(0)
	db.DB.Model(&Account{}).Where("user_id = ?", c.MustGet("user").(user.User).ID).Count(&count)
	return count
}

// delete 删除账户
func delete(c *gin.Context) error {
	account := c.MustGet("account").(*Account)

	return db.DB.Delete(account).Error
}

func edit(c *gin.Context, form accountEditForm) (*AccountDto, error) {
	account := c.MustGet("account").(*Account)

	if form.Name != "" {
		account.Name = form.Name
	}
	if form.Description != "" {
		account.Description = form.Description
	}
	if form.ExpectTotalMoney > 0 {
		account.ExpectTotalMoney = form.ExpectTotalMoney
	}
	if form.PerPartMoney > 0 {
		account.PerPartMoney = form.PerPartMoney
	}
	if form.ExpectRateOfReturn > 0 {
		account.ExpectRateOfReturn = form.ExpectRateOfReturn
	}

	if err := db.DB.Save(account).Error; err != nil {
		return nil, err
	}

	// 将 account 转换为 AccountDto
	accountDto := account.ToDto()

	return accountDto, nil
}

func GetByIdInternal(id uint) (*Account, error) {
	account := &Account{}
	if err := db.DB.First(account, id).Error; err != nil {
		return nil, ErrAccountNotFound
	}

	return account, nil
}

func GetByIdWithUidInternal(id uint, uid uint) (*Account, error) {
	account, err := GetByIdInternal(id)
	if err != nil {
		return nil, err
	}
	if account.UserId != uid {
		return nil, ErrAccountNotFound
	}

	return account, nil
}

type accountBindBreedForm struct {
	Id uint `form:"id" json:"id" binding:"required" label:"品种id"`
}

func bindBreed(ctx *gin.Context) (*AccountBreedDto, error) {
	var form accountBindBreedForm
	if err := ctx.ShouldBind(&form); err != nil {
		return nil, err
	}
	breedId := form.Id
	b, err := breed.GetByIdWithUidInternal(breedId, ctx.MustGet("user").(user.User).ID)
	if err != nil {
		return nil, err
	}
	accountBreed := AccountBreed{
		Account: *ctx.MustGet("account").(*Account),
		Breed:   *b,
		BreedId: b.ID,
		Model: model.Model{
			CreatedAt: carbon.Now().Timestamp(),
			UpdatedAt: carbon.Now().Timestamp(),
		},
	}
	// 得到breed 后就可以进行绑定了
	account := ctx.MustGet("account").(*Account)
	if err := db.DB.Model(account).Association("Breeds").Append([]AccountBreed{accountBreed}); err != nil {
		return nil, err
	}

	return accountBreed.ToDto(), nil
}

type accountAddBreedBuyItemForm struct {
	Id       uint    `form:"id" json:"id" binding:"required" label:"品种id"`
	CreateAt int64   `form:"create_at" json:"create_it" binding:"required" label:"买入时间"`
	Cost     float64 `form:"cost" json:"cost" binding:"required" label:"购买单价"`
	Count    int64   `form:"count" json:"count" binding:"required" label:"购买份数"`
}

func addBreedBuytItem(ctx *gin.Context) (*AccountBreedDto, error) {
	return nil, nil
}
