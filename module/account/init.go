package account

import "github.com/crazyhl/yzyx-materials/internal/db"

func AutoMigrate() {
	db.DB.AutoMigrate(&Account{}, &AccountBreed{}, &BuyBreedItem{})
}
