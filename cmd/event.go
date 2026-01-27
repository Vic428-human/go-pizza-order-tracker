package main

import "github.com/gin-gonic/gin"

// 通知處理器內部使用
func (h *Handler) ServeNotification(c *gin.Context) {

	topic := c.Query("orderId")
	client := make(chan string)
	// 透過 Handler 存取 Notification
	h.Notification.Subscribe(topic, client)

	// SSE/WebSocket 長連接...
}
