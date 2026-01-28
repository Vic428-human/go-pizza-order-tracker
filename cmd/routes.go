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
	router.Use(sessions.Sessions("pizza-tracker", store))

	// ====== TMPL 版本 ======
	router.GET("/", h.ServeNewOrderForm)
	router.POST("/new-order", h.HandleNewOrderPost)
	router.GET("/customer/:id", h.serveCustomer)
	router.GET("/notifications", h.ServeNotification)

	// ====== React 版本 ======

	router.GET("/login", h.HandleLoginGet)
	router.POST("/login", h.HandleLoginPost)
	router.POST("/logout", h.HandleLogoutPost)

	admin := router.Group("/admin")
	admin.Use(h.AuthMiddleware())
	{
		admin.GET("", h.ServeAdminDashboard)
		//透過 id 來更新訂單
		admin.POST("/order/:id/update", h.handleOrderPut)
		// 透過 id 來刪除訂單
		admin.POST("/order/:id/delete", h.handleOrderDelete)
		// 顧客送出添加訂單publish後，通知 admin
		admin.GET("/notifications", h.adminNotificationHandler)
	}

	// ====== TMPL 版本 ====== 把數據直接交付給 templtate
	router.Static("/static", "./templates/static")

	// ====== React API 版本 ======
	api := router.Group("/api")
	{
		adminApi := api.Group("/admin")
		adminApi.Use(h.AuthMiddleware())
		{
			adminApi.GET("/dashboard", h.GetAdminDashboardJSON)
		}
	}

	// ====== React 前端頁面 ======
	router.Static("/react", "./frontend/dist")

}
