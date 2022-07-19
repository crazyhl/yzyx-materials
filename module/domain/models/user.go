package models

import (
	"github.com/crazyhl/yzyx-materials/internal/model"
	"github.com/crazyhl/yzyx-materials/module/domain/dtos"
	"github.com/golang-jwt/jwt"
)

type User struct {
	model.Model
	Username string `gorm:"type:varchar(32);not null;unique"` // 用户名
	Password string `gorm:"type:varchar(128);not null"`       // 密码
}

//自定义Claims
type UserJwtClaims struct {
	ID       uint   `json:"id"`
	UserName string `json:"username"`
	jwt.StandardClaims
}

func (u User) ToDto() *dtos.UserDto {
	return &dtos.UserDto{
		ID:       u.ID,
		Username: u.Username,
	}
}
