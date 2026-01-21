package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// https://zhuanlan.zhihu.com/p/30184285330
func (h *Handler) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := GetSession(c, "userID")

		// 1. 判斷是否有 userID，空字串代表根本沒登入，自然 session 也不會存在。
		if userID == "" {
			// http.StatusSeeOther (303) => 重導向，指示客戶端在重導向後，必須改用 GET 方法來存取新的 URL，表單送出後，重導向到結果頁面（避免使用者重新整理時重複送出 POST）。
			// http.StatusTemporaryRedirect (307) => 適合用在「暫時性」的導向，例如：使用者未登入，先導到 /login，但登入後還要回到原本的請求。保留原本的請求語義（尤其是 POST 或 PUT），避免被強制轉成 GET。
			c.Redirect(http.StatusSeeOther, "/login")
			// Abort() => 只會停止後續的 middleware/handler 執行。
			// AbortWithStatusJSON() => 同樣會停止後續的 middleware/handler。 同時設定回應：直接設定 HTTP 狀態碼，並輸出 JSON 內容。
			c.AbortWithStatusJSON(401, gin.H{"error": "unauthorized"})
			return
		}

		// 2. 查詢資料庫確認使用者是否存在 => 使用者曾經登入過， session 還在，但資料庫裡已經沒有這個使用者
		user, err := h.users.GetUserByID(userID)
		if err != nil {
			// 資料庫運行異常
			c.Redirect(http.StatusSeeOther, "/login")
			c.AbortWithStatusJSON(500, gin.H{"error": "資料庫錯誤"})
			return
		}
		if user == nil {
			// 資料庫運行正常，但使用者不存在
			c.Redirect(http.StatusSeeOther, "/login")
			c.AbortWithStatusJSON(404, gin.H{"error": "資料庫裡已經沒有這個使用者"})
			return
		}

		// 3. 驗證成功 => 執行 route group /admin 中的 handler
		c.Next()
	}
}
