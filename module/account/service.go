package account

import (
	"fmt"

	"github.com/crazyhl/yzyx-materials/internal/db"
	"github.com/crazyhl/yzyx-materials/module/domain/dtos"
	"github.com/crazyhl/yzyx-materials/module/domain/models"
	"github.com/gin-gonic/gin"
)

// addAccount 添加账户
func add(form accountAddForm) (*dtos.AccountDto, error) {
	account := &models.Account{
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
func list(c *gin.Context) []*dtos.AccountDto {
	accounts := []*models.Account{}
	accountDtos := []*dtos.AccountDto{}
	db.DB.Scopes(db.Paginate(c)).Where("user_id = ?", c.MustGet("user").(models.User).ID).Order("id desc").Find(&accounts)
	for _, account := range accounts {
		accountDtos = append(accountDtos, account.ToDto())
	}
	return accountDtos
}

func getCount(c *gin.Context) int64 {
	count := int64(0)
	db.DB.Model(&models.Account{}).Where("user_id = ?", c.MustGet("user").(models.User).ID).Count(&count)
	return count
}

// delete 删除账户
func delete(c *gin.Context) error {
	account := c.MustGet("account").(*models.Account)

	return db.DB.Delete(account).Error
}

func edit(c *gin.Context, form accountEditForm) (*dtos.AccountDto, error) {
	account := c.MustGet("account").(*models.Account)

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

func GetByIdInternal(id uint) (*models.Account, error) {
	account := &models.Account{}
	if err := db.DB.First(account, id).Error; err != nil {
		return nil, ErrAccountNotFound
	}

	return account, nil
}

func GetByIdWithUidInternal(id uint, uid uint) (*models.Account, error) {
	account, err := GetByIdInternal(id)
	if err != nil {
		return nil, err
	}
	if account.UserId != uid {
		return nil, ErrAccountNotFound
	}

	return account, nil
}

func UpdateAccountProfit(breedId uint) {
	fmt.Println("breedId", breedId)
	accountIds := make([]uint, 0)
	db.DB.Table("account_breeds").Where("breed_id = ?", breedId).Pluck("account_id", &accountIds)
	fmt.Println("accountIds", accountIds)
	var accounts []models.Account
	db.DB.Preload("Breeds.Breed").Where("id in ?", accountIds).Find(&accounts)
	for _, account := range accounts {
		totalCost := float64(0)
		totalProfit := float64(0)
		for _, b := range account.Breeds {
			totalCost += b.TotalCost //  加入每个账户的总投入
			// 计算利润 品种的总市值 - 品种总投入
			totalProfit += (b.Breed.NetValue*float64(b.TotalCount) - b.TotalCost)
		}
		account.TotalMoney = totalCost
		account.ProfitAmount = totalProfit
		account.RateOfReturn = totalProfit / totalCost * 100
		db.DB.Save(&account)
	}
}
