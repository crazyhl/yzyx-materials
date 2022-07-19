package models

import "github.com/crazyhl/yzyx-materials/internal/db"

func AutoMigrate() {
	db.DB.AutoMigrate(&User{}, &Breed{}, &Account{}, &AccountBreed{}, &BuyBreedItem{})
}
