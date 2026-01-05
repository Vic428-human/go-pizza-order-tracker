package main

import (
	"github.com/gin-gonic/gin"
)

type OrderFormData struct {
	PizzaTypes []string
	PizzaSizes []string
}

// dive 是 go-playground/validator 提供的特殊標籤。
// 它用於啟用對 slice/array/map 內部元素的遞歸驗證。
// 若結構體中包含嵌套的切片或數組，且需要驗證其內部字段，必須加上 dive，否則只會驗證外層容器本身（如長度），不會驗證內部元素的字段。

type OrderRequest struct {
	Name         string   `json:"name" binding:"required,min=2,max=100"`
	Phone        string   `json:"phone" binding:"required, min=10,max=20"`
	Address      string   `json:"address" binding:"required,min=5,max=200"`
	Size         string   `json:"size" binding:"required,dive,valid_pizza_size"`
	PizzaTypes   []string `json:"pizzaTypes" binding:"required,dive,valid_pizza_type"`
	Instructions string   `json:"instructions" binding:"max=500"`
}

func (h *Handler) ServiceNewOrderList(c *gin.Context) {
	c.HTML()
}