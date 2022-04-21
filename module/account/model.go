package account

import (
	"github.com/crazyhl/yzyx-materials/internal/model"
	"github.com/crazyhl/yzyx-materials/module/user"
)

// 记账账户
type Account struct {
	model.Model
	Name             string    `gorm:"type:varchar(100);not null"`             // 账户名称
	Description      string    `gorm:"type:varchar(255);default: '';not null"` // 账户描述，描述一下投资目的
	UserId           uint      `gorm:"not null"`                               // 用户ID
	User             user.User // 账户所属用户，外键
	ExpectTotalMoney float64   `gorm:"type:decimal(20,2);default: 0;not null"` // 预期投入总金额
	TotalMoney       float64   `gorm:"type:decimal(20,2);default: 0;not null"` // 已经投入总金额
	PerPartMoney     float64   `gorm:"type:decimal(20,2);default: 0;not null"` // 每份金额
}
