// Package 引入
package main

import "pizza-tracker-go/internal/models"

// Handler 結構體
// 用來管理應用的主要邏輯層跟服務層
// orders 指向 OrderModel 指標，用來處理有關訂單的資料操作，像是建立/查詢
type Handler struct {
	// 可以透過 import 引用其他 models 
	orders *models.OrderModel
	// 也可以直接調用 cmd 目錄下的其他功能模組
}


// 1. 單元測試時，可以注入假的 MockOrderModel，輕鬆模擬資料庫行為
// 2. 未來換成 Redis 或其他資料庫，只需要修改 NewHandler 的呼叫處理即可
// 3. 配置靈活，可切換 dev 跟 prod環境
func NewHandler(dbModel *models.DBModel) *Handler {
	return &Handler{orders: &dbModel.Order}
}




