package account

import (
	"github.com/crazyhl/yzyx-materials/internal/db"
	"github.com/crazyhl/yzyx-materials/module/user"
	"github.com/gin-gonic/gin"
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
	db.DB.Scopes(db.Paginate(c)).Where("user_id = ?", c.MustGet("user").(user.User).ID).Find(&accounts)
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
func delete(c *gin.Context, id uint) error {
	account, err := GetByIdWithUidInternal(id, c.MustGet("user").(user.User).ID)
	if err != nil {
		return err
	}

	return db.DB.Delete(account).Error
}

func edit(c *gin.Context, id uint, form accountEditForm) error {
	account, err := GetByIdWithUidInternal(id, c.MustGet("user").(user.User).ID)
	if err != nil {
		return err
	}

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

	return db.DB.Save(account).Error
}

func GetByIdInternal(id uint) (*Account, error) {
	account := &Account{}
	if err := db.DB.First(account, id).Error; err != nil {
		return nil, err
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
