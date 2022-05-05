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

// delete DeleteBreed 删除购买品种
func delete(ctx *gin.Context, id uint) error {
	breed := &Breed{}
	db.DB.First(breed, id)
	if breed.Account.User.ID != ctx.MustGet("user").(user.User).ID {
		return ErrBreedNotYourAccount
	}

	return db.DB.Delete(&Breed{}, "id = ?", id).Error
}

func getBreedByIdInternal(id uint) (*Breed, error) {
	breed := &Breed{}
	if err := db.DB.Preload("Account.User").First(breed, id).Error; err != nil {
		return nil, err
	}

	return breed, nil
}

func getBreedByIdWithUidInternal(id, uid uint) (*Breed, error) {
	breed, err := getBreedByIdInternal(id)
	if err != nil {
		return nil, err
	}

	if breed.Account.User.ID != uid {
		return nil, ErrBreedNotYourAccount
	}

	return breed, nil
}

func addBreedBuyItem(ctx *gin.Context, form AddBreedItemForm) error {
	breed, err := getBreedByIdInternal(form.BreedID)
	if err != nil {
		return err
	}
	totalMoney := 0.0
	if form.Type == 1 {
		totalMoney = form.Cost*float64(form.TotalPart) + form.Commission
	} else {
		totalMoney = form.Cost*float64(form.TotalPart) - form.Commission

	}

	accPerPartMoneyTotalPart := 0.0
	if breed.Account.PerPartMoney > 0 {
		accPerPartMoneyTotalPart = totalMoney / breed.Account.PerPartMoney
	}

	breedBuyItem := &BreedBuyItem{
		Breed:                        *breed,
		Cost:                         form.Cost,
		TotalPart:                    form.TotalPart,
		TotalMoney:                   totalMoney,
		Commission:                   form.Commission,
		AccountPerPartMoneyTotalPart: accPerPartMoneyTotalPart,
		Type:                         form.Type,
	}

	if err := db.DB.Create(breedBuyItem).Error; err != nil {
		return err
	}
	// 插入够买记录后，更新 breed 相关字段 以及 account 相关字段
	// 不采用每次加减的方案，而是采用统计的方式
	updateBreedStatistics(breedBuyItem.Breed)
	return nil
}

type breedStatisticsResult struct {
	TotalPart      int
	TotalMoney     float64
	AccountPerPart float64
}

func updateBreedStatistics(b Breed) {
	buyStatResult := &breedStatisticsResult{}
	soldStatResult := &breedStatisticsResult{}
	db.DB.Model(&BreedBuyItem{}).Where("breed_id = ?", b.ID).Where("type = ?", 1).
		Select("sum(total_part) as total_part, sum(total_money) as total_money, sum(account_per_part_money_total_part) as account_per_part_money_total_part").
		First(buyStatResult)
	db.DB.Model(&BreedBuyItem{}).Where("breed_id = ?", b.ID).Where("type = ?", 2).
		Select("sum(total_part) as total_part, sum(total_money) as total_money, sum(account_per_part_money_total_part) as account_per_part_money_total_part").
		First(soldStatResult)
	b.TotalPart = buyStatResult.TotalPart - soldStatResult.TotalPart
	b.TotalMoney = buyStatResult.TotalMoney - soldStatResult.TotalMoney
	b.AccountPerPartMoneyTotalPart = buyStatResult.AccountPerPart - soldStatResult.AccountPerPart
	b.Cost = b.TotalMoney / float64(b.TotalPart)
	b.PercentForAccountExpectTotalMoney = b.TotalMoney / float64(b.Account.ExpectTotalMoney)
	b.PercentForAccountTotalMoney = b.TotalMoney / float64(b.Account.TotalMoney)
	db.DB.Save(&b)
}
