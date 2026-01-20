package main

import "github.com/gin-gonic/gin"

/*
第二個參數放的是 處理該路由的 handler function，發送請求時要執行的邏輯。
r.GET("/ping", func(c *gin.Context) { c.JSON(200, gin.H{"message": "pong"}) })
*/
func setupRoutes(router *gin.Engine, h *Handler) {
	router.GET("/", h.ServeNewOrderForm)            // 查看訂單
	router.POST("/new-order", h.HandleNewOrderPost) // 創建新的訂單
	router.GET("/customer/:id", h.serveCustomer)    // 查看顧客訂單內容
	// 基本实现：使用 router.Static 方法来设置静态文件目录，例如 router.Static("/static", "./public")，这将把 public 目录下的所有文件映射到 /static 路径。
	router.Static("/static", "./templates/static")
}
