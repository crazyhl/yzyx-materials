package account

import (
	"github.com/crazyhl/yzyx-materials/internal/model"
	"github.com/crazyhl/yzyx-materials/module/domain/models"
)

// 记账账户
type Account struct {
	model.Model
	Name               string      `gorm:"type:varchar(100);not null"`             // 账户名称
	Description        string      `gorm:"type:varchar(255);default: '';not null"` // 账户描述，描述一下投资目的
	UserId             uint        `gorm:"not null"`                               // 用户ID
	User               models.User // 账户所属用户，外键
	ExpectTotalMoney   float64     `gorm:"type:decimal(20,4);default: 0;not null"` // 预期投入总金额
	ExpectRateOfReturn uint8       `gorm:"type:tinyint;default: 0;not null"`       // 预期收益率
	RateOfReturn       float64     `gorm:"type:decimal(20,4);default: 0;not null"` // 实际收益率
	TotalMoney         float64     `gorm:"type:decimal(20,4);default: 0;not null"` // 已经投入总金额
	PerPartMoney       float64     `gorm:"type:decimal(20,4);default: 0;not null"` // 每份金额
	ProfitAmount       float64     `gorm:"type:decimal(20,4);default: 0;not null"` // 收益总金额
	Breeds             []*AccountBreed
}

func (a Account) ToDto() *AccountDto {
	breedDtos := make([]*AccountBreedDto, 0)
	for _, breed := range a.Breeds {
		breedDtos = append(breedDtos, breed.ToDto())
	}

	return &AccountDto{
		Dto: model.Dto{
			ID:        a.ID,
			CreatedAt: a.CreatedAt,
			UpdatedAt: a.UpdatedAt,
		},
		Name:               a.Name,
		Desc:               a.Description,
		Total:              a.TotalMoney,
		ExpectTotal:        a.ExpectTotalMoney,
		ExpectRateOfReturn: a.ExpectRateOfReturn,
		PerPart:            a.PerPartMoney,
		RateOfReturn:       a.RateOfReturn,
		ProfitAmount:       a.ProfitAmount,
		Breeds:             breedDtos,
	}
}

// AccountBreed 账户品种模型
type AccountBreed struct {
	model.Model
	AccountId                uint         `gorm:"not null"`
	Account                  Account      // 账户品种所属的账户
	BreedId                  uint         `gorm:"not null"`
	Breed                    models.Breed // 账户品种所属的品种
	TotalAccountPerPartCount float64      `gorm:"type:decimal(20,4);not null;default:0;"`  // 对应账户设置每份金额计算出来的总份数
	Cost                     float64      `gorm:"type:decimal(20,4);not null;default:0;"`  // 成本
	TotalCount               int64        `gorm:"not null;default:0;"`                     // 总份数
	TotalCost                float64      `gorm:"type:decimal(20,4);not null; default:0;"` // 总成本
	// 上面三个字段跟品种区别的就是这几个统计只属于该账户品种的
}

func (b *AccountBreed) ToDto() *AccountBreedDto {
	return &AccountBreedDto{
		Dto: model.Dto{
			ID:        b.ID,
			CreatedAt: b.CreatedAt,
			UpdatedAt: b.UpdatedAt,
		},
		Account:                  *b.Account.ToDto(),
		Breed:                    *b.Breed.ToDto(),
		Cost:                     b.Cost,
		TotalCount:               b.TotalCount,
		TotalCost:                b.TotalCost,
		TotalAccountPerPartCount: b.TotalAccountPerPartCount,
	}
}

// 账户购买品种记录
type BuyBreedItem struct {
	model.Model
	AccountId uint `gorm:"not null;"`
	Account   Account
	BreedId   uint `gorm:"not null;"`
	Breed     models.Breed
	Cost      float64 `gorm:"decimal(20,4);not null;default:0;"` // 成本
	Count     int64   `gorm:"not null;default:0"`                // 购买份数，如果是卖出则是负数
	TotalCost float64 `gorm:"decimal(20,4);not null;default:0;"` // 总成本
}
