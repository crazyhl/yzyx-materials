package breed

import (
	"github.com/crazyhl/yzyx-materials/internal/model"
	"github.com/crazyhl/yzyx-materials/module/domain/models"
)

// 品种 Model
type Breed struct {
	model.Model
	UserId        uint        `gorm:"not null;uniqueIndex: uk_uid_code;"` // 用户ID
	User          models.User // 账户所属用户，外键
	Code          string      `gorm:"type:varchar(32);not null;uniqueIndex: uk_uid_code;"` // 编码
	Name          string      `gorm:"type:varchar(32);not null;"`                          // 名称
	NetValue      float64     `gorm:"type:decimal(20,4);not null;default:0;"`              // 净值
	Cost          float64     `gorm:"type:decimal(20,4);not null;default:0;"`              // 成本
	TotalCount    int64       `gorm:"not null;default:0;"`                                 // 总份数
	TotalCost     float64     `gorm:"type:decimal(20,4);not null; default:0;"`             // 总成本
	TotalNetValue float64     `gorm:"type:decimal(20,4);not null; default:0"`              // 总净值
}

func (breed *Breed) ToDto() *BreedDto {
	// 将 account 转换为 AccountDto
	breedTto := &BreedDto{
		Code:          breed.Code,
		Name:          breed.Name,
		NetValue:      breed.NetValue,
		Cost:          breed.Cost,
		TotalCount:    breed.TotalCount,
		TotalCost:     breed.TotalCost,
		TotalNetValue: breed.TotalNetValue,
		Dto: model.Dto{
			ID:        breed.ID,
			CreatedAt: breed.CreatedAt,
			UpdatedAt: breed.UpdatedAt,
		},
	}

	return breedTto
}
