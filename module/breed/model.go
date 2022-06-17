package breed

import (
	"github.com/crazyhl/yzyx-materials/internal/model"
	"github.com/crazyhl/yzyx-materials/module/user"
)

// 品种 Model
type Breed struct {
	model.Model
	UserId   uint      `gorm:"not null"` // 用户ID
	User     user.User // 账户所属用户，外键
	Code     string    `gorm:"type:varchar(32);not null;unique"` // 编码
	Name     string    `gorm:"type:varchar(32);not null;unique"` // 名称
	NetValue float64   `gorm:"type:decimal(20,4);not null"`      // 净值
	Cost     float64   `gorm:"type:decimal(20,4);not null"`      // 成本
}

func (breed Breed) ToDto() *BreedDto {
	// 将 account 转换为 AccountDto
	breedTto := &BreedDto{
		Code:     breed.Code,
		Name:     breed.Name,
		NetValue: breed.NetValue,
		Cost:     breed.Cost,
	}

	return breedTto
}
