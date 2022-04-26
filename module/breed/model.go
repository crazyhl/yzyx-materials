package breed

import (
	"github.com/crazyhl/yzyx-materials/internal/model"
	"github.com/crazyhl/yzyx-materials/module/account"
)

// 品种
type Breed struct {
	model.Model
	AccountId                         uint `gorm:"not null;uniqueIndex:uk_acc_code;"` // 账户ID"`
	Account                           account.Account
	Code                              string  `gorm:"type:varchar(20);not null;uniqueIndex:uk_acc_code;"` // 品种编码
	Name                              string  `gorm:"type:varchar(20);not null;"`                         // 品种名称
	Cost                              float64 `gorm:"type:decimal(10,2);not null;"`                       // 成本
	Price                             float64 `gorm:"type:decimal(10,2);not null;"`                       // 当前价格 净值
	TotalPart                         int     `gorm:"type:int(11);not null;"`                             // 总份数
	TotalMoney                        float64 `gorm:"type:decimal(10,2);not null;"`                       // 总金额
	AccountPerPartMoneyTotalPart      int     `gorm:"type:int(11);not null;"`                             // 账户每份投入总份数
	TotalPrice                        float64 `gorm:"type:decimal(10,2);not null;"`                       // 总价格 持仓金额
	PercentForAccountExpectTotalMoney float64 `gorm:"type:decimal(10,2);not null;"`                       // 账户预计投入金额占比
	PercentForAccountTotalMoney       float64 `gorm:"type:decimal(10,2);not null;"`                       // 账户投入金额占比
}

type BreedBuyItem struct {
	model.Model
	BreedId                      uint `gorm:"not null;"` // 品种ID
	Breed                        Breed
	Cost                         float64 `gorm:"type:decimal(10,2);not null;"` // 单价成本
	TotalPart                    int     `gorm:"type:int(11);not null;"`       // 总份数
	TotalMoney                   float64 `gorm:"type:decimal(10,2);not null;"` // 总金额
	AccountPerPartMoneyTotalPart int     `gorm:"type:int(11);not null;"`       // 账户每份投入总份数
}
