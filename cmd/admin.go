package main

import (
	"fmt"
	"log"
	"net/http"
	"pizza-tracker-go/internal/models"

	"github.com/gin-gonic/gin"
)

type AdminOrderData struct {
	Orders   []models.Order
	Statuses []string
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
	var form struct {
		Account  string `form:"account" binding:"required,min=3,max=50"`
		Password string `form:"password" binding:"required,min=6"`
	}

	// 先判斷規則方面的錯誤
	if err := c.ShouldBind(&form); err != nil {
		c.HTML(http.StatusOK, "login.tmpl", LoginData{Error: "Invalid input: " + err.Error()})
		return
	}

	user, err := h.users.AuthenticateUser(form.Account, form.Password)

	//規則正確，但登入資訊錯誤
	if err != nil {
		c.HTML(http.StatusOK, "login.tmpl", LoginData{
			Error: err.Error(), // 使用實際的錯誤訊息
		})
		return
	}
	// 需要改成字串，因為等下操作 DB 的 GetUserByID 是拿 string去搜尋
	SetSession(c, "userID", fmt.Sprintf("%v", user.ID))
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

// 顯示所有訂單、先前存在session的username、所有訂單狀態
// orders => 所有訂單，每個訂單的實際進度條狀態，也就是當前狀態處在哪個階段
// Status => 需要把所有狀態傳進去是因為要做下拉選單，所以要知道總共有哪些狀態可以供選擇
func (h *Handler) ServeAdminDashboard(c *gin.Context) {
	orders, err := h.orders.GetAllOrders()
	if err != nil {
		log.Printf("獲取訂單資訊失敗!!!: %v", err)
		c.String(http.StatusInternalServerError, "獲取訂單資訊失敗!!!")
	}
	username := GetSession(c, "username")

	log.Printf("===>當前登入帳號: %s", username)
	c.HTML(http.StatusOK, "admin.tmpl", AdminOrderData{
		Orders:   orders,
		Statuses: models.OrderStatues,
		Username: username,
	})
}

// 更新特定訂單
func (h *Handler) handleOrderPut(c *gin.Context) {
	orderID := c.Param("id")
	newStatus := c.PostForm("status")

	// 更新狀態失敗
	if err := h.orders.UpdateOrderStatus(orderID, newStatus); err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	topic := "order:" + orderID
	h.notificationManager.Publish(topic, "訂單狀態已更新成 : "+newStatus)

	// 更新狀態成功
	c.Redirect(http.StatusSeeOther, "/admin")

}

// 刪除特定訂單
func (h *Handler) handleOrderDelete(c *gin.Context) {
	orderID := c.Param("id")
	if err := h.orders.DeleteOrder(orderID); err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	c.Redirect(http.StatusSeeOther, "/admin")
}
