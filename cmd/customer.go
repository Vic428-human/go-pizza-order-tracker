package main

import (
	"log/slog"
	"net/http"
	"os"
	"pizza-tracker-go/internal/models"

	"github.com/gin-gonic/gin"
)

type CustomerData struct {
	Title string
	Order models.Order
	Statuses []string
}

type OrderFormData struct { // 定義 從 models 取得披薩種類與尺寸的資料 的結構體
	PizzaTypes []string
	PizzaSizes []string
}

// dive 是 go-playground/validator 提供的特殊標籤，它用於啟用對 slice/array/map 內部元素的遞歸驗證，若結構體中包含嵌套的切片或數組，且需要驗證其內部字段，必須加上 dive，否則只會驗證外層容器本身（如長度），不會驗證內部元素的字段。
type OrderReuqest struct {
	Name         string   `form:"name" binding:"required,min=2,max=100"`
	Phone        string   `form:"phone" binding:"required,min=10,max=20"`
	Address      string   `form:"address" binding:"required,min=5,max=200"`
	Sizes        []string `form:"size" binding:"required,min=1,dive,valid_pizza_size"`
	PizzaTypes   []string `form:"pizza" binding:"required,min=1,dive,valid_pizza_type"`
	Instructions []string `form:"instructions" binding:"max=200"`
}

//  tmpl 前端模板
func (h *Handler) ServeNewOrderForm(c *gin.Context) { // ServeNewOrderForm 屬於 Handler 結構體的方法，用來處理 HTTP 請求。
	// 回傳一個 HTML 頁面
	c.HTML(http.StatusOK, "order.tmpl", OrderFormData{
		PizzaTypes: models.PizzaTypes, // 把 資料包裝成 OrderFormData結構體，然後提供給模板
		PizzaSizes: models.PizzaSizes,
	})
}

// Undefined validation function 'min' on field 'Phone' => 這錯誤跟 binding 的寫法錯誤有關
func (h *Handler) HandleNewOrderPost(c *gin.Context) {
	var form OrderReuqest

	if err := c.ShouldBind(&form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	/* 效果:
	[]models.OrderItem{
		{
			Size: "Large",
			Pizza: "Margherita",
			Instructions: "",
		},
		{
			Size: "Medium",
			Pizza: "Pepperoni",
			Instructions: "",
		},
	}
	
	*/
	// 組合訂單明細 : 準備一個清單，裝每一個pizza訂單項目
	orderItems := make([]models.OrderItem, len(form.Sizes))
	for i := range orderItems { // 把 表單的資料，一筆一筆的轉乘 OrderItem struct，將結果塞進 orderItems slice中
		orderItems[i] = models.OrderItem{ // 用意: 把訂單項目，變成有意義的物件，而不是零散的slice，也方便後續處理
			Size:         form.Sizes[i],
			Pizza:        form.PizzaTypes[i],
			Instructions: form.Instructions[i],
		}
	}
	
	order := models.Order{
		Status:       models.OrderStatues[0],
		CustomerName: form.Name,
		Phone:        form.Phone,
		Address:      form.Address,
		Items:        orderItems,
	}

	// 當前 func 已經跟 Handler 結構體綁定，可以直接透過 h.orders 呼叫 OrderModel 的方法
	if err := h.orders.CreateOrder(&order); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create order"})
		// 情況1:會多了 time / level等欄位說明
		// time=2026-01-08T00:23:00.000+08:00 level=ERROR msg="Failed to create order" error="some error message"
		// err.Error() 會返回具體的錯誤訊息字串放在 error 欄位
		logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
		logger.Error("處理請求失敗", "error", err.Error())  // 或 slog.Any("error", err)

		// 寫法2 傳統單行，會少了欄位說明: 
		// slog.Error("Failed to create order", "error", err)

		c.String(http.StatusInternalServerError, "Something went wrong")
		return
	}
	slog.Info("Order created", "orderId", order.ID, "customer", order.CustomerName)
	
	// 請求的資源可用，並且應該獲取 https://blog.csdn.net/weixin_42073635/article/details/143805554 
	// c.Redirect(statusCode, location)
	c.Redirect(http.StatusSeeOther, "/customer/" + order.ID)
	
}


func (h *Handler) serveCustomer(c *gin.Context) {

	// /customer/:id ( 假設:id前端忘了提供)
	orderID := c.Param("id")
	if orderID == "" {
		// 後端回傳
		c.String(http.StatusBadRequest, "請前端提供訂單ID")
	}
	// ( 假設:id 前端有提供，但可能提供的是錯誤的，資料庫找不到讓前端知道)
	order, err := h.orders.GetOrder(orderID)
	if err != nil {
		c.String(http.StatusNotFound, "資料庫中查不到該筆訂單")
		return
	}

	// 如果資料庫有訂單，以 tmpl 呈現給前端
	c.HTML(http.StatusOK, "customer.tmpl", CustomerData{
		Title: "Pizza Tracker - 訂單詳情"+ order.ID,
		Order: *order,
		Statuses: models.OrderStatues,
	})
}