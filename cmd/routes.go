package main

import "github.com/gin-gonic/gin"

/*
第二個參數放的是 處理該路由的 handler function，發送請求時要執行的邏輯。
r.GET("/ping", func(c *gin.Context) { c.JSON(200, gin.H{"message": "pong"}) })
*/
func setupRoutes(router *gin.Engine, h *Handler) {
	router.GET("/", h.ServeNewOrderForm)
}