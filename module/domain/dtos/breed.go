package dtos

import "github.com/crazyhl/yzyx-materials/internal/model"

type BreedDto struct {
	model.Dto
	Code          string  `json:"code"`            // 编码
	Name          string  `json:"name"`            // 名称
	NetValue      float64 `json:"net_value"`       // 净值
	Cost          float64 `json:"cost"`            // 成本
	TotalCount    int64   `json:"total_count"`     // 总份数
	TotalCost     float64 `json:"total_cost"`      // 总成本
	TotalNetValue float64 `json:"total_net_value"` // 总净值
}
