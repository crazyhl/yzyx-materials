package account

type AccountDto struct {
	ID                 uint    `json:"id"`                    // 账户ID
	Name               string  `json:"name"`                  // 账户名
	Desc               string  `json:"desc"`                  // 账户描述
	Total              float64 `json:"total"`                 // 账户总额
	ExpectTotal        float64 `json:"expect_total"`          // 预期总投入
	ExpectRateOfReturn uint8   `json:"expect_rate_of_return"` // 预期收益率
	RateOfReturn       float64 `json:"rate_of_return"`        // 当前投入部分收益率
	PerPart            float64 `json:"per_part"`              // 每份金额
	Created            int64   `json:"created"`               // 创建时间
	Updated            int64   `json:"updated"`               // 更新时间
	ProfitAmount       float64 `json:"profit_amount"`         // 收益总额
}
