package accountbreed

import (
	"github.com/crazyhl/yzyx-materials/internal/db"
	"github.com/crazyhl/yzyx-materials/internal/model"
	"github.com/crazyhl/yzyx-materials/module/breed"
	"github.com/crazyhl/yzyx-materials/module/domain/dtos"
	"github.com/crazyhl/yzyx-materials/module/domain/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-module/carbon/v2"
	"gorm.io/gorm/clause"
)

type accountAddBreedBuyItemForm struct {
	Id       uint    `form:"id" json:"id" binding:"required" label:"品种id"`
	CreateAt int64   `form:"create_at" json:"create_at" binding:"required" label:"买入时间"`
	Cost     float64 `form:"cost" json:"cost" binding:"required" label:"购买单价"`
	Count    int64   `form:"count" json:"count" binding:"required" label:"购买份数"`
	Fee      float64 `form:"fee" json:"fee" binding:"required" label:"手续费"`
}

type AccountBreedStatisticsResult struct {
	TotalCount int64
	TotalCost  float64
}

// 账户添加购买记录
func addBreedBuytItem(ctx *gin.Context) (*dtos.AccountBreedDto, error) {
	var form accountAddBreedBuyItemForm
	if err := ctx.ShouldBind(&form); err != nil {
		return nil, err
	}
	breedId := form.Id
	b, err := breed.GetByIdWithUidInternal(breedId, ctx.MustGet("user").(models.User).ID)
	if err != nil {
		return nil, err
	}
	account := ctx.MustGet("account").(*models.Account)
	buyItem := models.BuyBreedItem{
		Account: *account,
		Breed:   *b,
		Model: model.Model{
			CreatedAt: form.CreateAt,
		},
		Cost:      form.Cost,
		Count:     form.Count,
		Fee:       form.Fee,
		TotalCost: form.Cost*float64(form.Count) + form.Fee,
	}
	if err := db.DB.Create(&buyItem).Error; err != nil {
		return nil, err
	}
	afterChangeBuyItem(buyItem)
	breed := models.AccountBreed{}
	db.DB.Preload(clause.Associations).Where("account_id = ?", buyItem.Account.ID).Where("breed_id = ?", buyItem.Breed.ID).First(&breed)
	return (&breed).ToDto(), nil
}

// 账户添加购买记录
func updateBreedBuytItem(ctx *gin.Context) (*dtos.BuyBreedItemDto, error) {
	var form accountAddBreedBuyItemForm
	if err := ctx.ShouldBind(&form); err != nil {
		return nil, err
	}
	buyItem := models.BuyBreedItem{}
	db.DB.Preload(clause.Associations).Where("id = ?", form.Id).First(&buyItem)
	uid := ctx.MustGet("user").(models.User).ID
	if buyItem.Account.UserId != uid || buyItem.Breed.UserId != uid {
		return nil, ErrBuyItemNotFound
	}

	buyItem.Cost = form.Cost
	buyItem.Count = form.Count
	buyItem.Fee = form.Fee
	buyItem.TotalCost = form.Cost*float64(form.Count) + form.Fee
	if err := db.DB.Save(&buyItem).Error; err != nil {
		return nil, err
	}
	afterChangeBuyItem(buyItem)
	return buyItem.ToDto(), nil
}

func afterChangeBuyItem(buyItem models.BuyBreedItem) {
	// 添加成功后，更新 Breed, AccountBreed,以及 Account 三个表的数据
	// 获取 accountBreed
	breed := models.AccountBreed{}
	db.DB.Preload(clause.Associations).Where("account_id = ?", buyItem.Account.ID).Where("breed_id = ?", buyItem.Breed.ID).First(&breed)
	// 根据accountId 和 breedId 进行数据统计 更新 account_breed
	accountBreedStatisticsResult := AccountBreedStatisticsResult{}
	db.DB.Model(&models.BuyBreedItem{}).Where("account_id = ?", buyItem.Account.ID).Where("breed_id = ?", buyItem.Breed.ID).Select("sum(count) as total_count, sum(total_cost) as total_cost").First(&accountBreedStatisticsResult)
	breed.TotalCost = accountBreedStatisticsResult.TotalCost
	breed.TotalCount = accountBreedStatisticsResult.TotalCount
	breed.Cost = accountBreedStatisticsResult.TotalCost / float64(accountBreedStatisticsResult.TotalCount)
	if breed.Account.PerPartMoney > 0 {
		breed.TotalAccountPerPartCount = accountBreedStatisticsResult.TotalCost / float64(breed.Account.PerPartMoney)
	}
	db.DB.Save(&breed)
	// 根据 breed 进行统计 更新 breed
	db.DB.Model(&models.BuyBreedItem{}).Where("breed_id = ?", buyItem.Breed.ID).Select("sum(count) as total_count, sum(total_cost) as total_cost").First(&accountBreedStatisticsResult)
	buyItem.Breed.TotalCost = accountBreedStatisticsResult.TotalCost
	buyItem.Breed.TotalCount = accountBreedStatisticsResult.TotalCount
	buyItem.Breed.Cost = accountBreedStatisticsResult.TotalCost / float64(accountBreedStatisticsResult.TotalCount)
	db.DB.Save(buyItem.Breed)
	// 根据 account 进行统计，更新 account
	account := &models.Account{}
	db.DB.Preload("Breeds.Breed").First(account, buyItem.Account.ID)
	totalCost := float64(0)
	totalProfit := float64(0)
	for _, b := range account.Breeds {
		totalCost += b.TotalCost //  加入每个账户的总投入
		// 计算利润 品种的总市值 - 品种总投入
		totalProfit += b.Breed.NetValue*float64(b.TotalCount) - b.TotalCost
	}
	account.TotalMoney = totalCost
	account.ProfitAmount = totalProfit
	account.RateOfReturn = totalProfit / totalCost * 100
	db.DB.Save(account)
}

type accountBindBreedForm struct {
	Id uint `form:"id" json:"id" binding:"required" label:"品种id"`
}

func bindBreed(ctx *gin.Context) (*dtos.AccountBreedDto, error) {
	var form accountBindBreedForm
	if err := ctx.ShouldBind(&form); err != nil {
		return nil, err
	}
	breedId := form.Id
	b, err := breed.GetByIdWithUidInternal(breedId, ctx.MustGet("user").(models.User).ID)
	if err != nil {
		return nil, err
	}
	accountBreed := models.AccountBreed{
		Account: *ctx.MustGet("account").(*models.Account),
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

// buyItemList 账户的某个品种购买列表
func buyItemList(ctx *gin.Context) []dtos.BuyBreedItemDto {
	// 根据 account_id 和 breed_id 获取品种购买列表就ok了
	buyItems := make([]models.BuyBreedItem, 0)
	db.DB.Scopes(db.Paginate(ctx)).
		Where("account_id = ?", ctx.Param("id")).
		Where("breed_id = ?", ctx.Param("breedId")).
		Order("created_at DESC").Find(&buyItems)
	buyItemDtos := make([]dtos.BuyBreedItemDto, 0)
	for _, item := range buyItems {
		buyItemDtos = append(buyItemDtos, *item.ToDto())
	}
	return buyItemDtos
}

// getBuyItemCount 获取账户品种购买记录数量
func getBuyItemCount(ctx *gin.Context) int64 {
	count := int64(0)
	db.DB.Model(&models.BuyBreedItem{}).
		Where("account_id = ?", ctx.Param("id")).
		Where("breed_id = ?", ctx.Param("breedId")).
		Count(&count)
	return count
}
