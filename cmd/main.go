package main

// := 表示 同時聲明跟初始化變量

import (
	"fmt"
	"pizza-tracker-go/internal/models" // 導入db模塊
	"pizza-tracker-go/mathutil"        // 導入加法函數
)

func main() {
	// := 複製來自 pizza-tracker-go/internal/models 裡的db數據，函數內的修改不影響原先的值
	
	dbModel,err := models.InitDB("test.db")
    sum := mathutil.Add(3, 5)
    fmt.Println("結果:", sum)

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(dbModel)

}
