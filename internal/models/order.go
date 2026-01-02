package models

import (
	"time"

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

/*  假數據參考
[
  {
    "id": "ORD20260102001", //primaryKey , ID
    "status": "pending",
    "customerName": "張小明",
    "phone": "0912345678",
    "address": "台北市信義區信義路五段7號15樓",
    "items": [
      {
        "id": 1, //primaryKey
        "order_id": "ORD20260102001", // 外鍵，指向 Order.ID
        "product": "iPhone 16 Pro 256GB 鈦金屬灰",
        "quantity": 1,
        "unitPrice": 35900
      },
      {
        "id": 2, //primaryKey
        "order_id": "ORD20260102001", // 外鍵，指向 Order.ID
        "product": "AirPods Pro 2 USB-C版",
        "quantity": 2,
        "unitPrice": 7990
      },
    ],
	createdAt: "2026-01-02T00:00:00.000Z"
  },
 
]
*/ 

// 這裡大寫的ID / Items 其實是欄位名稱
type Order struct {
	ID string `gorm:"primaryKey;uniqueIndex;size:14" json:"id"`
	Status string `gorm:"not null" json:"status"`
	CustomerName string `gorm:"not null" json:"customerName"`
	Phone string `gorm:"not null" json:"phone"`
	Address string `gorm:"not null" json:"address"`
	// 一對多關聯，在 OrderItem 裡有訂單ID (OrderID)，指向的是 Order 裡的 ID (Order.ID)
	Items []OrderItem `gorm:"foreignKey:OrderID" json:"items"`
	CreatedAt  time.Time `json:"createdAt"`
}

type OrderItem struct {
	ID string `gorm:"primaryKey;size:14" json:"id"`
	OrderID string `gorm:"size:14;index" json:"order_id"` // 外鍵，指向 Order.ID
}
