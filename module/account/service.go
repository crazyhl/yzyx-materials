package account

import (
	"github.com/crazyhl/yzyx-materials/internal/db"
	"github.com/crazyhl/yzyx-materials/internal/model"
	"github.com/crazyhl/yzyx-materials/module/breed"
	"github.com/crazyhl/yzyx-materials/module/user"
	"github.com/gin-gonic/gin"
	"github.com/golang-module/carbon/v2"
	"gorm.io/gorm/clause"
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
		Model: model.Model{
			CreatedAt: carbon.Now().Timestamp(),
			UpdatedAt: carbon.Now().Timestamp(),
		},
	}
	// 得到breed 后就可以进行绑定了
	if err := db.DB.Create(&accountBreed).Error; err != nil {
		return nil, err
	}

	return accountBreed.ToDto(), nil
}

type accountAddBreedBuyItemForm struct {
	Id       uint    `form:"id" json:"id" binding:"required" label:"品种id"`
	CreateAt int64   `form:"create_at" json:"create_at" binding:"required" label:"买入时间"`
	Cost     float64 `form:"cost" json:"cost" binding:"required" label:"购买单价"`
	Count    int64   `form:"count" json:"count" binding:"required" label:"购买份数"`
}

type AccountBreedStatisticsResult struct {
	TotalCount int64
	TotalCost  float64
}

// 账户添加购买记录
func addBreedBuytItem(ctx *gin.Context) (*AccountBreedDto, error) {
	var form accountAddBreedBuyItemForm
	if err := ctx.ShouldBind(&form); err != nil {
		return nil, err
	}
	breedId := form.Id
	b, err := breed.GetByIdWithUidInternal(breedId, ctx.MustGet("user").(user.User).ID)
	if err != nil {
		return nil, err
	}
	account := ctx.MustGet("account").(*Account)
	buyItem := BuyBreedItem{
		Account: *account,
		Breed:   *b,
		Model: model.Model{
			CreatedAt: form.CreateAt,
		},
		Cost:      form.Cost,
		Count:     form.Count,
		TotalCost: form.Cost * float64(form.Count),
	}
	if err := db.DB.Create(&buyItem).Error; err != nil {
		return nil, err
	}
	// 添加成功后，更新 Breed, AccountBreed,以及 Account 三个表的数据
	// 获取 accountBreed
	breed := AccountBreed{}
	db.DB.Preload(clause.Associations).Where("account_id = ?", buyItem.Account.ID).Where("breed_id = ?", buyItem.Breed.ID).First(&breed)
	// 根据accountId 和 breedId 进行数据统计 更新 account_breed
	accountBreedStatisticsResult := AccountBreedStatisticsResult{}
	db.DB.Model(&BuyBreedItem{}).Where("account_id = ?", buyItem.Account.ID).Where("breed_id = ?", buyItem.Breed.ID).Select("sum(count) as total_count, sum(total_cost) as total_cost").First(&accountBreedStatisticsResult)
	breed.TotalCost = accountBreedStatisticsResult.TotalCost
	breed.TotalCount = accountBreedStatisticsResult.TotalCount
	breed.Cost = accountBreedStatisticsResult.TotalCost / float64(accountBreedStatisticsResult.TotalCount)
	if breed.Account.PerPartMoney > 0 {
		breed.TotalAccountPerPartCount = accountBreedStatisticsResult.TotalCost / float64(breed.Account.PerPartMoney)
	}
	db.DB.Save(&breed)
	// 根据 breed 进行统计 更新 breed
	db.DB.Model(&BuyBreedItem{}).Where("breed_id = ?", buyItem.Breed.ID).Select("sum(count) as total_count, sum(total_cost) as total_cost").First(&accountBreedStatisticsResult)
	b.TotalCost = accountBreedStatisticsResult.TotalCost
	b.TotalCount = accountBreedStatisticsResult.TotalCount
	b.Cost = accountBreedStatisticsResult.TotalCost / float64(accountBreedStatisticsResult.TotalCount)
	db.DB.Save(b)
	// 根据 account 进行统计，更新 account
	allAccount := &Account{}
	db.DB.Preload(clause.Associations).First(allAccount, account.ID)
	totalCost := float64(0)
	totalProfit := float64(0)
	for _, b := range allAccount.Breeds {
		totalCost += b.TotalCost //  加入每个账户的总投入
		// 计算利润 品种的总市值 - 品种总投入
		totalProfit += b.Breed.NetValue*float64(b.TotalCount) - b.TotalCost
	}
	allAccount.TotalMoney = totalCost
	allAccount.ProfitAmount = totalProfit
	allAccount.RateOfReturn = totalProfit / totalCost * 100
	db.DB.Save(allAccount)
	db.DB.Preload(clause.Associations).Where("account_id = ?", buyItem.Account.ID).Where("breed_id = ?", buyItem.Breed.ID).First(&breed)
	return (&breed).ToDto(), nil
}
