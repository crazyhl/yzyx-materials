package breed

import (
	"github.com/crazyhl/yzyx-materials/internal/bus"
	"github.com/crazyhl/yzyx-materials/internal/db"
	"github.com/crazyhl/yzyx-materials/module/domain/dtos"
	"github.com/crazyhl/yzyx-materials/module/domain/models"
	"github.com/gin-gonic/gin"
)

// add 添加账户
func add(form addForm) (*dtos.BreedDto, error) {
	breed := &models.Breed{
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
func edit(ctx *gin.Context, form editForm) (*dtos.BreedDto, error) {
	breed := ctx.MustGet("breed").(*models.Breed)

	breed.Code = form.Code
	breed.Name = form.Name
	if form.NetValue > 0 {
		breed.NetValue = form.NetValue
		breed.TotalNetValue = float64(breed.TotalCount) * form.NetValue
	}

	if form.Cost > 0 {
		breed.Cost = form.Cost
	}

	if err := db.DB.Save(breed).Error; err != nil {
		return nil, err
	}

	// 将 account 转换为 AccountDto
	breedTto := breed.ToDto()

	return breedTto, nil
}

// delete 删除账户
func delete(ctx *gin.Context) error {
	breed := ctx.MustGet("breed").(*models.Breed)

	return db.DB.Delete(breed).Error
}

// listAccounts 获取账户列表
func list(c *gin.Context) []*dtos.BreedDto {
	breeds := []*models.Breed{}
	breedDtos := []*dtos.BreedDto{}
	query := db.DB.Scopes(db.Paginate(c)).Where("user_id = ?", c.MustGet("user").(models.User).ID).Order("id desc")
	filter := c.Query("filter")
	if filter != "" {
		query.Where(db.DB.Where("code like ?", "%"+filter+"%").Or("name like ?", "%"+filter+"%"))
	}
	query.Find(&breeds)
	for _, breed := range breeds {
		breedDtos = append(breedDtos, breed.ToDto())
	}
	return breedDtos
}

func getCount(c *gin.Context) int64 {
	count := int64(0)
	db.DB.Model(&models.Breed{}).Where("user_id = ?", c.MustGet("user").(models.User).ID).Count(&count)
	return count
}

// updateNetValue 更新净值
func updateNetValue(ctx *gin.Context, netValue float64) (*dtos.BreedDto, error) {
	breed := ctx.MustGet("breed").(*models.Breed)

	breed.NetValue = netValue
	breed.TotalNetValue = float64(breed.TotalCount) * netValue

	if err := db.DB.Save(breed).Error; err != nil {
		return nil, err
	}
	// 获取这个用户绑定了品种的账户，对有这个品种的账户进行数据更新
	// account.UpdateAccountProfit(breed.ID)
	bus.Bus.Publish("account:updateProfit", breed.ID)
	// 将 account 转换为 AccountDto
	breedTto := breed.ToDto()

	return breedTto, nil
}

// getByIdInternal 获取账户
func GetByIdInternal(id uint) (*models.Breed, error) {
	breed := &models.Breed{}
	if err := db.DB.First(breed, id).Error; err != nil {
		return nil, ErrBreedNotFound
	}

	return breed, nil
}

// getByIdWithUidInternal 获取账户并校验 uid
func GetByIdWithUidInternal(id uint, uid uint) (*models.Breed, error) {
	breed, err := GetByIdInternal(id)
	if err != nil {
		return nil, err
	}
	if breed.UserId != uid {
		return nil, ErrBreedNotFound
	}

	return breed, nil
}
