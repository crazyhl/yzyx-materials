package user

import (
	"github.com/crazyhl/yzyx-materials/internal/model"
)

type User struct {
	model.Model
	Username string `gorm:"type:varchar(32);not null;unique"` // 用户名
	Password string `gorm:"type:varchar(32);not null"`        // 密码
}
