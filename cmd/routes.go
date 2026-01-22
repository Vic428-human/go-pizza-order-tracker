package main

import (
	"github.com/gin-contrib/sessions"
	gormsessions "github.com/gin-contrib/sessions/gorm"
	"github.com/gin-gonic/gin"
)

/*
第二個參數放的是 處理該路由的 handler function，發送請求時要執行的邏輯。
r.GET("/ping", func(c *gin.Context) { c.JSON(200, gin.H{"message": "pong"}) })
*/
func setupRoutes(router *gin.Engine, h *Handler, store gormsessions.Store) {

	// 是瀏覽器 cookie 的名稱（例如 Set-Cookie: pizza-tracker=...
	// store：指定 session 資料儲存在哪裡（這裡是 SQLite）
	router.Use(sessions.Sessions("pizza-tracker", store))

	// 不用登入就可以訪問的endpoint
	router.GET("/", h.ServeNewOrderForm)            // 查看訂單
	router.POST("/new-order", h.HandleNewOrderPost) // 創建新的訂單
	router.GET("/customer/:id", h.serveCustomer)    // 查看顧客訂單內容

	router.GET("/login", h.HandleLoginGet)
	router.POST("/login", h.HandleLoginPost)
	router.POST("/logout", h.HandleLogoutPost)

	// 有登入後才能訪問
	admin := router.Group("/admin")
	admin.Use(h.AuthMiddleware())
	{
		// admin.GET("", h.ServeAdminDashboard)       // 顯示後台首頁
		// admin.POST("/order/:id/update", h.HandleOrderPut) // 更新訂單
	}

	// 把 templates/static 目錄下的所有文件映射到 /static 路徑。
	router.Static("/static", "./templates/static")

}
