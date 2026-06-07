package validator

import (
	"fmt"
	"yanblog/utils/errmsg"
	"github.com/go-playground/locales/zh_Hans_CN"
	unTrans "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	zhTrans "github.com/go-playground/validator/v10/translations/zh"
	"reflect"
)

// Validate 验证数据结构并返回翻译后的错误信息
// 参数: data - 需要验证的数据结构
// 返回: string - 错误信息, int - 状态码
func Validate(data interface{}) (string, int) {
	validate := validator.New() //  创建一个新的验证器
	uni := unTrans.New(zh_Hans_CN.New()) //  创建一个通用翻译器
	trans, _ := uni.GetTranslator("zh_Hans_CN") //  获取指定语言的翻译器

	err := zhTrans.RegisterDefaultTranslations(validate, trans) //  注册默认的翻译器
	if err != nil {
		fmt.Println("err:", err)
	}
	validate.RegisterTagNameFunc(func(field reflect.StructField) string {
		label := field.Tag.Get("label")
		return label
	})

	err = validate.Struct(data) //  验证数据
	if err != nil {
		for _, v := range err.(validator.ValidationErrors) {
			return v.Translate(trans), errmsg.ERROR
		}
	}
	return "", errmsg.SUCCESS
}
