// Package 引入
package main

import "pizza-tracker-go/internal/models"

// 用來管理應用的主要邏輯層跟服務層
type Handler struct {
	orders              *models.OrderModel
	users               *models.UserModel
	notificationManager *NotificationManager
}

// 1. 單元測試時，可以注入假的 MockOrderModel，輕鬆模擬資料庫行為
// 2. 未來換成 Redis 或其他資料庫，只需要修改 NewHandler 的呼叫處理即可
// 3. 配置靈活，可切換 dev 跟 prod環境
func NewHandler(dbModel *models.DBModel) *Handler {
	return &Handler{
		orders:              &dbModel.Order,
		users:               &dbModel.User,
		notificationManager: NewNotificationManager(),
	}
}
