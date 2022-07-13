package breed

import "github.com/crazyhl/yzyx-materials/internal/model"

type BreedDto struct {
	model.Dto
	Code          string  `json:"code"`                      // 编码
	Name          string  `json:"name"`                      // 名称
	NetValue      float64 `json:"net_value,omitempty"`       // 净值
	Cost          float64 `json:"cost,omitempty"`            // 成本
	TotalCount    int64   `json:"total_count,omitempty"`     // 总份数
	TotalCost     float64 `json:"total_cost,omitempty"`      // 总成本
	TotalNetValue float64 `json:"total_net_value,omitempty"` // 总净值
}
