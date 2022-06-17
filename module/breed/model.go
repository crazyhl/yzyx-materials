package breed

import "github.com/crazyhl/yzyx-materials/internal/model"

// 品种 Model
type Breed struct {
	model.Model
	Code     string  `gorm:"type:varchar(32);not null;unique"` // 编码
	Name     string  `gorm:"type:varchar(32);not null;unique"` // 名称
	NetValue float64 `gorm:"type:decimal(20,4);not null"`      // 净值
	Cost     float64 `gorm:"type:decimal(20,4);not null"`      // 成本
}
