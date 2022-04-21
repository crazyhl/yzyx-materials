package model

// model 自己的 baseModel 实现
type Model struct {
	ID        uint  `gorm:"primarykey"`                          // 自增主键
	CreatedAt int64 `gorm:"autoCreateTime;not null;default: 0;"` // 创建时间
	UpdatedAt int64 `gorm:"autoUpdateTime;not null;default: 0;"` // 更新时间
}
