package breed

import (
	"github.com/crazyhl/yzyx-materials/internal/db"
)

// add 添加账户
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

// edit 编辑账户
func edit(form editForm, id, uid uint) (*BreedDto, error) {
	breed, err := getByIdWithUidInternal(id, uid)
	if err != nil {
		return nil, err
	}

	breed.Code = form.Code
	breed.Name = form.Name
	breed.NetValue = form.NetValue
	breed.Cost = form.Cost

	if err := db.DB.Save(breed).Error; err != nil {
		return nil, err
	}

	// 将 account 转换为 AccountDto
	breedTto := breed.ToDto()

	return breedTto, nil
}

// getByIdInternal 获取账户
func getByIdInternal(id uint) (*Breed, error) {
	breed := &Breed{}
	if err := db.DB.First(breed, id).Error; err != nil {
		return nil, ErrBreedNotFound
	}

	return breed, nil
}

// getByIdWithUidInternal 获取账户并校验 uid
func getByIdWithUidInternal(id uint, uid uint) (*Breed, error) {
	breed, err := getByIdInternal(id)
	if err != nil {
		return nil, err
	}
	if breed.UserId != uid {
		return nil, ErrBreedNotFound
	}

	return breed, nil
}
