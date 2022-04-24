package account

import (
	"github.com/crazyhl/yzyx-materials/internal/db"
	"github.com/crazyhl/yzyx-materials/module/user"
	"github.com/gin-gonic/gin"
)

// addAccount 添加账户
func add(form accountAddForm) (*AccountDto, error) {
	account := &Account{
		Name:             form.Name,
		Description:      form.Description,
		User:             form.User,
		ExpectTotalMoney: form.ExpectTotalMoney,
		PerPartMoney:     form.PerPartMoney,
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
