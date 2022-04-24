package account

import (
	"github.com/crazyhl/yzyx-materials/internal"
)

func add(form accountAddForm) (*AccountDto, error) {
	account := &Account{
		Name:             form.Name,
		Description:      form.Description,
		User:             form.User,
		ExpectTotalMoney: form.ExpectTotalMoney,
		PerPartMoney:     form.PerPartMoney,
	}
	if err := internal.DB.Create(account).Error; err != nil {
		return nil, err
	}

	// 将 account 转换为 AccountDto
	accountDto := &AccountDto{
		ID:      account.ID,
		Name:    account.Name,
		Desc:    account.Description,
		Total:   account.TotalMoney,
		Expect:  account.ExpectTotalMoney,
		PerPart: account.PerPartMoney,
		Created: account.CreatedAt,
		Updated: account.UpdatedAt,
	}

	return accountDto, nil
}
