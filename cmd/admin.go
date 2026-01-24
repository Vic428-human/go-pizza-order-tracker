package main

import (
	"net/http"
	"pizza-tracker-go/internal/models"

	"github.com/gin-gonic/gin"
)

type AdminData struct {
	Orders   []models.Order
	Status   []string
	Username string
}

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
		Password string `json:"password" binding:"required,min=6"`
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

	SetSession(c, "userID", user.ID)
	SetSession(c, "username", user.Username)

	// 登入成功，存session跟導轉路徑
	// SetSession(c, "userID", user.ID)
	c.Redirect(http.StatusSeeOther, "/admin")
}

func (h *Handler) HandleLogoutPost(c *gin.Context) {

	if err := ClearAllSession(c); err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	c.Redirect(http.StatusSeeOther, "/login")
}

// 需要顯示所有訂單、先前存在session的username、所有訂單狀態
// orders => 所有訂單，每個訂單的實際進度條狀態，也就是當前狀態處在哪個階段
// Status => 需要把所有狀態傳進去是因為要做下拉選單，所以要知道總共有哪些狀態可以供選擇
func (h *Handler) ServeAdminDashboard(c *gin.Context) {
	orders, _ := h.orders.GetAllOrders()
	username := GetSession(c, "username")
	c.HTML(http.StatusOK, "admin.tmpl", AdminData{Orders: orders, Status: models.OrderStatues, Username: username})
}
