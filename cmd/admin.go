package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// 處理登入邏輯 => session不存在時，導轉去login頁面，此時要把錯誤訊息顯示再登入頁面
type LoginData struct {
	Error string
}

// c *gin.Context => 代表一次 HTTP 請求與回應的上下文。透過它可以讀取請求、回傳資料。
func (h *Handler) HandleLoginGet(c *gin.Context) {
	c.HTML(http.StatusOK, "login.tmpl", LoginData{}) // LoginData{} → 傳入模板的資料
}
