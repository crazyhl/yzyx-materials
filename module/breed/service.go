package breed

import (
	"github.com/crazyhl/yzyx-materials/internal/db"
	"github.com/crazyhl/yzyx-materials/module/user"
	"github.com/gin-gonic/gin"
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

// delete 删除账户
func delete(uid, id uint) error {
	breed, err := getByIdWithUidInternal(id, uid)
	if err != nil {
		return err
	}

	return db.DB.Delete(breed).Error
}

// listAccounts 获取账户列表
func list(c *gin.Context) []*BreedDto {
	breeds := []*Breed{}
	breedDtos := []*BreedDto{}
	db.DB.Scopes(db.Paginate(c)).Where("user_id = ?", c.MustGet("user").(user.User).ID).Order("id desc").Find(&breeds)
	for _, breed := range breeds {
		breedDtos = append(breedDtos, breed.ToDto())
	}
	return breedDtos
}

func getCount(c *gin.Context) int64 {
	count := int64(0)
	db.DB.Model(&Breed{}).Where("user_id = ?", c.MustGet("user").(user.User).ID).Count(&count)
	return count
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
