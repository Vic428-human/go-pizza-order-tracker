package main

import (
	"net/http"
	"pizza-tracker-go/internal/models"

	"github.com/gin-gonic/gin"
)

type OrderFormData struct {
	PizzaTypes []string
	PizzaSizes []string
}

// dive 是 go-playground/validator 提供的特殊標籤，它用於啟用對 slice/array/map 內部元素的遞歸驗證，若結構體中包含嵌套的切片或數組，且需要驗證其內部字段，必須加上 dive，否則只會驗證外層容器本身（如長度），不會驗證內部元素的字段。
type OrderRequest struct {
	Name         string   `json:"name" binding:"required,min=2,max=100"`
	Phone        string   `json:"phone" binding:"required, min=10,max=20"`
	Address      string   `json:"address" binding:"required,min=5,max=200"`
	Sizes        []string `form:"size" binding:"required,min=1,dive,valid_pizza_size"`
	PizzaTypes   []string `form:"pizza" binding:"required,min=1,dive,valid_pizza_type"`
	Instructions []string `form:"instructions" binding:"max=200"`
}



func (h *Handler) ServiceNewOrderList(c *gin.Context) {
	c.HTML(http.StatusOK, "order.tmpl", OrderFormData{
		PizzaTypes: models.PizzaTypes,
		PizzaSizes: models.PizzaSizes,
	})
}

func (h *Handler) HandleNewOrderPost(c *gin.Context) {
	var form OrderRequest

	if err := c.ShouldBind(&form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 組合訂單明細
	orderItems := make([]models.OrderItem, len(form.Sizes))
	for i := range orderItems {
		orderItems[i] = models.OrderItem{
			Size:         form.Sizes[i],
			Pizza:        form.PizzaTypes[i],
			Instructions: form.Instructions[i],
		}
	}
	
	
}
