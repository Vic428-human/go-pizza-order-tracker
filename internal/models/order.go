package models

import (
	"time"

	"github.com/teris-io/shortid" // https://bbs.itying.com/topic/687b507f4715aa008848880f ex: iNove6iQ9J / NVDve6-9Q
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

// db 包裝成一個 struct ，方便在多處使用， 在文檔中會透過 db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{}) 方式創建 db 實例，但封裝成 OrderModel struct 也是一種方式
// 透過 DB 指向(*) *gorm.DB 這個實例，有了 gorm 實例就可以用它具備的 方法 Create 也是其中之一
type OrderModel struct {
	DB *gorm.DB
}

/*  假數據參考
[
  {
    "id": "ORD20260102001", //primaryKey ，ID，可以透過 shortid package去產出
    "status": "pending",
    "customerName": "張小明",
    "phone": "0912345678",
    "address": "台北市信義區信義路五段7號15樓",
    "items": [ items 裡的訂單明細，id是唯一的
      {
        "id": 1, //primaryKey ，ID，可以透過 shortid package去產出 
        "order_id": "ORD20260102001", // 外鍵，指向 Order.ID， order_id 有透過外鍵連結到 Order.ID，所以是繼承 Order的 ID，所以相同
        "product": "iPhone 16 Pro 256GB 鈦金屬灰",
        "quantity": 1,
        "unitPrice": 35900
      },
      {
        "id": 2, //primaryKey ，ID，可以透過 shortid package去產出
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
	ID string `gorm:"primaryKey;size:14" json:"id"`
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
	// index tag 用途: gorm 在執行 automigrate 或 creteTable 時，自動創建資料庫索引(非唯一)
	// 加入 index 索引優點:
	// 1. 查詢訂單(Order)的明細 SELECT * FROM order_items WHERE order_id = 'some_order_id';
	// 2. 可快速定位紀錄，避免全表搜索
	OrderID string `gorm:"size:14;index;not null" json:"order_id"` // 外鍵，指向 Order.ID
	Size string `gorm:"not null" json:"size"`
	Pizza string `gorm:"not null" json:"pizza"`
	Instructions string `json:"instructions"` 
}

// 在 db.Create() 操作之前，這些 hook 都會自動被呼叫 (hook ex: BeforeCreate / AfterCreate / CreateOrder 等等)
func (o *Order) BeforeCreate(tx *gorm.DB) error{
	if(o.ID == "" ){
		o.ID = shortid.MustGenerate() // 確保 ID 總是生成，若失敗則 panic 中止操作，避免無效記錄插入資料庫
	}
	return nil
}

func (oi *OrderItem) BeforeCreate(tx *gorm.DB) error{
	if(oi.ID == "" ){
		oi.ID = shortid.MustGenerate() // 確保 ID 總是生成，若失敗則 panic 中止操作，避免無效記錄插入資料庫
	}
	return nil
}

// 執行實際的 SQL INSERT 語句到資料庫（這才是真正的「CreateOrder」完成的部分），所以會在 BeforeCreate 之後發生
func (o *OrderModel) CreateOrder(order *Order) error{
	return o.DB.Create(order).Error
}