package main

import (
	"log/slog"
	"os"
	"pizza-tracker-go/internal/models"
)

func main() {

	cfg := loadConfig()

	// 預留客製化的設置，在 production 情況才觸發
	// 但還不確定甚麼情況需要特別設定
	opts := &slog.HandlerOptions{
		Level:     slog.LevelInfo, // production 通常使用 XX 级别
		AddSource: false,          // production 關閉源碼位置以提高性能
	}

	var handler slog.Handler
	env := os.Getenv("ENV") 
	if env == "" { 
		env = "development"
 	}

	if env == "development" { // 預設...
		handler = slog.NewTextHandler(os.Stdout, nil)
	} else {  // 其他環境變數的情況
		handler = slog.NewTextHandler(os.Stdout, opts)
	}
	slog.SetDefault(slog.New(handler))

	// 1. 先初始化DB，連接DB，接著才處理結構體可以使用tag規則
	// dbModel := &DBModel{ 
	// 	DB: db, // *gorm.DB 把sqlite的 db gorm物件覆寫
	// 	Order: OrderModel{DB: db}, // 把sqlite的 db gorm物件覆寫在 OrderModel裡的 DB
	// }
	dbModel, err := models.InitDB(cfg.DBPath) 
	
	if err != nil {
		slog.Error("資料庫初始化失敗", "error", err)
		os.Exit(1)
	}

	slog.Info("資料庫連接成功", "path", cfg.DBPath)
	// 處理結構體可以使用tag規則
	RegisterCustomValidators()

	// Handler 可以理解成綁定了資料庫跟對應的模組裡的方法
	h := NewHandler(dbModel) // 因為已經跟資料庫連接所以也綁定 Order 這個欄位
	// Order這個欄位對應的model結構體是 OrderModel，而OrderModel結構體綁定過的方法都可以跟著使用

	
}