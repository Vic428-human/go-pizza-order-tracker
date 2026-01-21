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

	// 處理登入邏輯 => session不存在時，導轉去login頁面，此時要把錯誤訊息顯示再登入頁面
	router.GET("/login", h.HandleLoginGet)

	// 有登入後才能訪問
	admin := router.Group("/admin")
	admin.Use(h.AuthMiddleware())

	// 把 templates/static 目錄下的所有文件映射到 /static 路徑。
	router.Static("/static", "./templates/static")

}
