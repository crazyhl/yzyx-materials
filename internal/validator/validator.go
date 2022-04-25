package validator

import (
	"github.com/gin-gonic/gin/binding"
	"github.com/gookit/validate"
	"github.com/gookit/validate/locales/zhcn"
)

// implements the binding.StructValidator
type customValidator struct{}

func (c *customValidator) ValidateStruct(ptr interface{}) error {
	v := validate.Struct(ptr)
	v.Validate() // do validating
	zhcn.Register(v)

	if v.Errors.Empty() {
		return nil
	}
	return v.Errors
}

func (c *customValidator) Engine() interface{} {

	return nil
}

// 初始化 validation
func init() {
	// 更改全局选项
	zhcn.RegisterGlobal()
	validate.Config(func(opt *validate.GlobalOption) {
		opt.ValidateTag = "binding"
	})
	binding.Validator = &customValidator{}
}
