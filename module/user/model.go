package user

import (
	"github.com/crazyhl/yzyx-materials/internal/model"
	"github.com/golang-jwt/jwt"
)

type User struct {
	model.Model
	Username string `gorm:"type:varchar(32);not null;unique"` // 用户名
	Password string `gorm:"type:varchar(128);not null"`       // 密码
}

//自定义Claims
type UserJwtClaims struct {
	ID       uint
	UserName string
	jwt.StandardClaims
}
