package model

// model 自己的 baseModel 实现
type Model struct {
	ID        uint  `gorm:"primarykey"`     // 自增主键
	CreatedAt int64 `gorm:"autoCreateTime"` // 创建时间
	UpdatedAt int64 `gorm:"autoUpdateTime"` // 更新时间
}
