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

func (h *Handler) HandleLoginPost(c *gin.Context) {
	type LoginRequest struct {
		Username string `json:"username" binding:"required,min=3,max=50"`
		Password string `json:"password" binding:"required,min=6"` // 至少6個字元，實際可再加強
	}

	var req LoginRequest
	if err := c.ShouldBind(&req); err != nil {
		c.HTML(http.StatusOK, "login.tmpl", LoginData{Error: "Invalid input: " + err.Error()})
		return
	}

	user, err := h.users.AuthenticateUser(req.Username, req.Password)
	// 登入失敗
	if err != nil {
		c.HTML(http.StatusOK, "login.tmpl", LoginData{
			Error: "Invalid credentials",
		})
		return
	}

	// 登入成功，存session跟導轉路徑
	// SetSession(c, "userID", user.ID)
	c.Redirect(http.StatusSeeOther, "/admin")
}
