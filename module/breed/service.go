package breed

import "github.com/crazyhl/yzyx-materials/internal/db"

func add(form addForm) (*BreedDto, error) {
	breed := &Breed{
		Code:     form.Code,
		Name:     form.Name,
		NetValue: form.NetValue,
		Cost:     form.Cost,
		User:     form.User,
	}
	if err := db.DB.Create(breed).Error; err != nil {
		return nil, err
	}

	// 将 account 转换为 AccountDto
	breedTto := breed.ToDto()

	return breedTto, nil
}
