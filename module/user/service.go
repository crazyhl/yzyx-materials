package user

import (
	"errors"
	"time"

	"github.com/crazyhl/yzyx-materials/internal"
	"github.com/golang-jwt/jwt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
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
var ErrGetJWTError = errors.New("登录失败") // 获取用户 jwt 失败

// login 登录
func login(username, password string) (*UserDto, error) {
	userDto := &UserDto{}
	user := &User{}
	internal.DB.Where("username = ?", username).First(&user)
	if user.ID > 0 {
		// 找到用户了
		// 验证密码
		err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
		if err != nil {
			return userDto, ErrUserNotFound
		}

		userDto.ID = user.ID
		userDto.Username = user.Username
		tokenStr, err := generateJWT(userDto)
		if err != nil {
			return userDto, ErrGetJWTError
		}
		userDto.Token = tokenStr

		return userDto, nil
	} else {
		// 没有找到用户
		return userDto, ErrUserNotFound
	}
}

// generateJWT 生成JWT
func generateJWT(user *UserDto) (string, error) {
	claims := &UserJwtClaims{
		ID:       user.ID,
		UserName: user.Username,
	}

	claims.ExpiresAt = time.Now().Add(time.Hour * 24 * 7).Unix()
	// Parse the token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte(viper.GetString("JWT_SECRET")))
	if err != nil {
		log.Error("Generate JWT err: ", err)
	}
	return tokenString, err
}
