package breed

import "github.com/crazyhl/yzyx-materials/internal/db"

// add AddBreed 添加购买品种
func add(form AddBreedForm) (*BreedDto, error) {
	breed := &Breed{
		Code:      form.Code,
		Name:      form.Name,
		AccountId: form.AccountId,
	}
	if err := db.DB.Create(breed).Error; err != nil {
		return nil, err
	}

	return breed.ToDto(), nil
}
