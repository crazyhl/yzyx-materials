package user

import (
	"github.com/crazyhl/yzyx-materials/internal/db"
	"github.com/golang-jwt/jwt"
	"github.com/golang-module/carbon/v2"
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
	result := db.DB.Create(user)

	if user.ID > 0 {
		return user, nil
	}

	log.Error("register user err: ", result.Error)
	return nil, result.Error
}

// login 登录
func login(username, password string) (*UserDto, error) {
	userDto := &UserDto{}
	user := &User{}
	db.DB.Where("username = ?", username).First(&user)
	if user.ID > 0 {
		// 找到用户了
		// 验证密码
		err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
		if err != nil {
			return userDto, ErrUserNotFound
		}

		tokenStr, err := GenerateJWT(*user)
		if err != nil {
			return userDto, ErrGetJWTError
		}
		userDto.ID = user.ID
		userDto.Username = user.Username
		userDto.Token = tokenStr

		return userDto, nil
	} else {
		// 没有找到用户
		return userDto, ErrUserNotFound
	}
}

// generateJWT 生成JWT
func GenerateJWT(user User) (string, error) {
	claims := &UserJwtClaims{
		ID:       user.ID,
		UserName: user.Username,
	}

	claims.ExpiresAt = carbon.Now().AddDays(7).Timestamp()
	// Parse the token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte(viper.GetString("JWT_SECRET")))
	if err != nil {
		log.Error("Generate JWT err: ", err)
	}
	return tokenString, err
}

// ParseJWT 转换解析JWT
func ParseJwt(authorization string) (*UserJwtClaims, error) {
	jwtStringHeader := authorization[0:6]
	if jwtStringHeader == "Bearer" {
		jwtStringBody := authorization[7:]
		token, err := jwt.ParseWithClaims(jwtStringBody, &UserJwtClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(viper.GetString("JWT_SECRET")), nil
		})

		if err != nil {
			log.Error("Parse JWT err: ", err)
			return nil, ErrorJWTInvalid
		}

		if claims, ok := token.Claims.(*UserJwtClaims); ok && token.Valid {
			// 验证时间
			return claims, nil
		} else {
			return nil, ErrorJWTInvalid
		}

	} else {
		return nil, ErrorJWTInvalid
	}
}

// GetByUid 根据uid获取用户
func GetByUid(id uint) (User, error) {
	user := &User{}
	result := db.DB.First(&user, id)
	if result.Error != nil {
		return *user, result.Error
	}
	return *user, nil
}
