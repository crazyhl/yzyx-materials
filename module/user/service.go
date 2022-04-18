package user

import (
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
