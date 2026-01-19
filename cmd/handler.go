// Package 引入
package main

import "pizza-tracker-go/internal/models"

// Handler 結構體
// 用來管理應用的主要邏輯層跟服務層
// orders 指向 OrderModel 指標，用來處理有關訂單的資料操作，像是建立/查詢
type Handler struct {
	// 可以透過 import 引用其他 models
	// 這裡的 orders/users就是資料表的名稱，像是INSERT INTO的時候抓的就是這裡的欄位
	orders *models.OrderModel
	users  *models.UserModel

	// 也可以直接調用 cmd 目錄下的其他功能模組
	// ...
}

// 1. 單元測試時，可以注入假的 MockOrderModel，輕鬆模擬資料庫行為
// 2. 未來換成 Redis 或其他資料庫，只需要修改 NewHandler 的呼叫處理即可
// 3. 配置靈活，可切換 dev 跟 prod環境
func NewHandler(dbModel *models.DBModel) *Handler {
	// dbModel 是一個指向 models.DBModel 的指標。
	// DBModel 結構體裡面有一個欄位 Order，型別應該是 models.OrderModel。
	// dbModel.Order 代表的是 取出 DBModel 裡的 Order 欄位。
	// & 是「取址運算子」，意思是「取得某個變數的記憶體位址」。
	return &Handler{orders: &dbModel.Order} // &dbModel.Order 則是 指向這個值的指標 (OrderModel)*。
	// 這樣 Handler 就能透過 orders 操作 OrderModel 的方法。
}
