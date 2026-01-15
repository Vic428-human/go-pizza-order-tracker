package main

import (
	"pizza-tracker-go/internal/models"
	"slices"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10" // 註冊自訂驗證規則，讓你在模型結構體的欄位標籤（struct tag）中可以使用這些規則來驗證輸入資料是否合法
)


func RegisterCustomValidators() {
    if v, ok := binding.Validator.Engine().(*validator.Validate); ok {

        // 註冊自定義校驗方法:
		// models.PizzaSizes => createSliceValidator 它接收你想要允許的值清單 ex: []string{"Small", "Medium", "Large",}
        if err := v.RegisterValidation("valid_pizza_size", createSliceValidator(models.PizzaSizes)); err != nil {
			panic(err)
        }

		if err := v.RegisterValidation("valid_pizza_type", createSliceValidator(models.PizzaTypes)); err != nil {
			panic(err)
		}
        return
    }
	panic("validator engine is not of type *validator.Validate")
}

// 自定義檢測規則
func createSliceValidator(allowed []string) validator.Func {
	// 任何符合 func(fl validator.FieldLevel) bool 簽名的函式，就是 validator.Func
	return func(fl validator.FieldLevel) bool {
		value := fl.Field().String() // 取得欄位字串值
		return slices.Contains(allowed, value)
	}
 
}
