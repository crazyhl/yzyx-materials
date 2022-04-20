package user

// UserDto 用户数据传输对象
type UserDto struct {
	ID       uint   `json:"id"`       // 用户ID
	Username string `json:"username"` // 用户名
	Token    string `json:"token"`    // JWT
}
