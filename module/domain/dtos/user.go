package dtos

type UserDto struct {
	ID       uint   `json:"id"`              // 用户ID
	Username string `json:"username"`        // 用户名
	Token    string `json:"token,omitempty"` // JWT
}
