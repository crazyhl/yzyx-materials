package model

import "gorm.io/gorm"

// model 自己的 baseModel 实现
type Model struct {
	ID        uint           `gorm:"primarykey"`      // 自增主键
	CreatedAt int64          `gorm:"autoCreateTime;"` // 创建时间
	UpdatedAt int64          `gorm:"autoUpdateTime;"` // 更新时间
	DeletedAt gorm.DeletedAt `gorm:"index"`           // 软删除标志
}

type Dto struct {
	ID        uint  `json:"id"`
	CreatedAt int64 `json:"create_at,omitempty"`
	UpdatedAt int64 `json:"update_at,omitempty"`
}
