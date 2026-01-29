package main

import (
	"io"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

// 通知顧客，當前這筆訂單狀態剛剛改變
func (h *Handler) notificationHandler(c *gin.Context) {
	orderId := c.Query("orderId")

	if orderId == "" {
		c.String(http.StatusBadRequest, "orderId is required")
		return
	}

	_, err := h.orders.GetOrder(orderId)
	if err != nil {
		c.String(http.StatusNotFound, "order not found")
		return
	}
	topic := "order:" + orderId
	client := make(chan string, 10)

	// 訂閱 order:XXX 這個頻道，之後 admin 更改訂單狀態時，都會把訊息發送給有訂閱這個訂單的 clients
	h.notificationManager.Subscribe(topic, client)

	// 不管這個函數從哪裡 return（400、404、或是正常跑完 SSE），進入 defer 的 code 一定會在函數結束前執行一次。
	defer func() {
		h.notificationManager.Unsubscribe(topic, client)
		slog.Info("Unsubscribed from topic", "orderId", topic)
	}()

	h.streamSSE(c, client)
}

// 通知 ADMIN 有新的訂單請求
func (h *Handler) StreamNewOrderNotifications(c *gin.Context) {
	topic := "admin:new_orders"

	client := make(chan string, 10)

	h.notificationManager.Subscribe(topic, client)

	// 不管這個函數從哪裡 return（400、404、或是正常跑完 SSE），進入 defer 的 code 一定會在函數結束前執行一次。
	defer func() {
		h.notificationManager.Unsubscribe(topic, client)
		slog.Info("Unsubscribed from topic", "orderId", topic)
	}()

	h.streamSSE(c, client)
}

// SSE/WebSocket 長連接...訂閱後的固定寫法
func (h *Handler) streamSSE(c *gin.Context, client chan string) {
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")

	c.Stream(func(w io.Writer) bool {
		if msg, ok := <-client; ok {
			c.SSEvent("message", msg)
			return true
		}
		return false
	})
}
