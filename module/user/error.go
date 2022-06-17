package user

import "errors"

var ErrorJWTInvalid = errors.New("验证失败，请重新登录") // 无效的 JWT
var ErrUserNotFound = errors.New("用户名或密码错误")
var ErrGetJWTError = errors.New("登录失败") // 获取用户 jwt 失败
var ErrJWTExpired = errors.New("登录已过期，请重新登录")
