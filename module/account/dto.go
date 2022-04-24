package account

type AccountDto struct {
	ID      uint    `json:"id"`       // 账户ID
	Name    string  `json:"name"`     // 账户名
	Desc    string  `json:"desc"`     // 账户描述
	Total   float64 `json:"total"`    // 账户总额
	Expect  float64 `json:"expect"`   // 账户预期总额
	PerPart float64 `json:"per_part"` // 每份金额
	Created int64   `json:"created"`  // 创建时间
	Updated int64   `json:"updated"`  // 更新时间
}
