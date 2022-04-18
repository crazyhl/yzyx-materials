package user

import (
	"errors"

	"github.com/crazyhl/yzyx-materials/internal"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

// Register 用户注册
func register(username, password string) (*User, error) {
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Error("generate password err: ", err)
		return nil, err
	}
	user := &User{
		Username: username,
		Password: string(hashPassword),
	}
	result := internal.DB.Create(user)

	if user.ID > 0 {
		return user, nil
	}

	log.Error("register user err: ", result.Error)
	return nil, result.Error
}

var ErrUserNotFound = errors.New("用户名或密码错误")

func login(username, password string) (*User, error) {
	user := &User{}
	internal.DB.Where("username = ?", username).First(&user)
	if user.ID > 0 {
		// 找到用户了
		// 验证密码
		err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
		if err != nil {
			return user, ErrUserNotFound
		}
		return user, nil
	} else {
		// 没有找到用户
		return user, ErrUserNotFound
	}
}
