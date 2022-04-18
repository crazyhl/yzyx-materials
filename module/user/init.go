package user

import "github.com/crazyhl/yzyx-materials/internal"

func AutoMigrate() {
	internal.DB.AutoMigrate(&User{})
}
