package breed

import (
	"github.com/crazyhl/yzyx-materials/internal/model"
	"github.com/crazyhl/yzyx-materials/module/account"
)

// 品种
type Breed struct {
	model.Model
	AccountId                         uint `gorm:"not null;"` // 账户ID"`
	Account                           account.Account
	Code                              string  `gorm:"type:varchar(20);"`            // 品种编码
	Name                              string  `gorm:"type:varchar(20);not null;"`   // 品种名称
	Cost                              float64 `gorm:"type:decimal(10,2);not null;"` // 成本
	Price                             float64 `gorm:"type:decimal(10,2);not null;"` // 当前价格 净值
	TotalPart                         int     `gorm:"type:int(11);not null;"`       // 总份数
	TotalMoney                        float64 `gorm:"type:decimal(10,2);not null;"` // 总金额
	AccountPerPartMoneyTotalPart      int     `gorm:"type:int(11);not null;"`       // 账户每份投入总份数
	TotalPrice                        float64 `gorm:"type:decimal(10,2);not null;"` // 总价格 持仓金额
	PercentForAccountExpectTotalMoney float64 `gorm:"type:decimal(10,2);not null;"` // 账户预计投入金额占比
	PercentForAccountTotalMoney       float64 `gorm:"type:decimal(10,2);not null;"` // 账户投入金额占比
}

type BreedDto struct {
	model.Dto
	Code                              string  `json:"code"`
	Name                              string  `json:"name"`
	Cost                              float64 `json:"cost"`
	Price                             float64 `json:"price"`
	TotalPart                         int     `json:"total_part"`
	TotalMoney                        float64 `json:"total_money"`
	AccountPerPartMoneyTotalPart      int     `json:"account_per_part_money_total_part"`
	TotalPrice                        float64 `json:"total_price"`
	PercentForAccountExpectTotalMoney float64 `json:"percent_for_account_expect_total_money"`
	PercentForAccountTotalMoney       float64 `json:"percent_for_account_total_money"`
}

func (b Breed) ToDto() *BreedDto {
	return &BreedDto{
		Dto: model.Dto{
			ID:        b.ID,
			CreatedAt: b.CreatedAt,
			UpdatedAt: b.UpdatedAt,
		},
		Code:                              b.Code,
		Name:                              b.Name,
		Cost:                              b.Cost,
		Price:                             b.Price,
		TotalPart:                         b.TotalPart,
		TotalMoney:                        b.TotalMoney,
		AccountPerPartMoneyTotalPart:      b.AccountPerPartMoneyTotalPart,
		TotalPrice:                        b.TotalPrice,
		PercentForAccountExpectTotalMoney: b.PercentForAccountExpectTotalMoney,
		PercentForAccountTotalMoney:       b.PercentForAccountTotalMoney,
	}
}

type BreedBuyItem struct {
	model.Model
	BreedId                      uint `gorm:"not null;"` // 品种ID
	Breed                        Breed
	Type                         uint8   `gorm:"not null;index;"`              // 类型, 1 买入 2 卖出
	Cost                         float64 `gorm:"type:decimal(10,2);not null;"` // 单价成本
	TotalPart                    uint    `gorm:"type:int(11);not null;"`       // 总份数
	Commission                   float64 `gorm:"type:decimal(10,2);not null;"` // 手续费
	TotalMoney                   float64 `gorm:"type:decimal(10,2);not null;"` // 总金额
	AccountPerPartMoneyTotalPart float64 `gorm:"type:decimal(10,2);not null;"` // 账户每份投入总份数
}
