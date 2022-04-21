package account

import "github.com/crazyhl/yzyx-materials/internal"

func add(form accountAddForm) (*Account, error) {
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
	return account, nil
}
