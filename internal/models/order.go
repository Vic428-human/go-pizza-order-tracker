package models

import (
	"gorm.io/gorm"
)


var (
	OrderStatues = []string{
		"Order placed",
		"Preparing",
		"Baking",
		"Quality Check",
		"Ready",
	}

	PizzaTypes = []string{
		"Cheese",
		"Pepperoni",
		"Veggie",
	}

	PizzaSizes = []string{
		"Small",
		"Medium",
		"Large",
	}
)

type OrderModel struct {
	DB *gorm.DB
}