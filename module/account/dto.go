package account

import (
	"github.com/crazyhl/yzyx-materials/internal/model"
	"github.com/crazyhl/yzyx-materials/module/breed"
)

type AccountDto struct {
	model.Dto
	Name               string            `json:"name"`                  // 账户名
	Desc               string            `json:"desc"`                  // 账户描述
	Total              float64           `json:"total"`                 // 账户总额
	ExpectTotal        float64           `json:"expect_total"`          // 预期总投入
	ExpectRateOfReturn uint8             `json:"expect_rate_of_return"` // 预期收益率
	RateOfReturn       float64           `json:"rate_of_return"`        // 当前投入部分收益率
	PerPart            float64           `json:"per_part"`              // 每份金额
	ProfitAmount       float64           `json:"profit_amount"`         // 收益总额
	Breeds             []*breed.BreedDto `json:"breeds,omitempty"`      // 账户绑定的品种
}
