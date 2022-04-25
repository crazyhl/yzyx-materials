package validator

import (
	"reflect"
	"strings"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	zh_translations "github.com/go-playground/validator/v10/translations/zh"
)

var trans ut.Translator

// 初始化 validation 翻译器
func init() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterTagNameFunc(func(field reflect.StructField) string {
			label := field.Tag.Get("label")
			if label == "" {
				return field.Name
			}
			return label
		})
		zh := zh.New()
		uni := ut.New(zh, zh)
		trans, _ = uni.GetTranslator("zh")
		zh_translations.RegisterDefaultTranslations(v, trans)
	}

}

func Translate(err error) string {
	errs := err.(validator.ValidationErrors)
	transErrors := errs.Translate(trans)
	errMsgs := make([]string, 0)
	for _, e := range transErrors {
		errMsgs = append(errMsgs, e)
	}

	return strings.Join(errMsgs, "\n")
}
