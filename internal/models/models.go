package models

// 為什麼 github.com/glebarez/sqlite 可以正常執行
// 這個套件是 純 Go 實作的 SQLite driver。
// 它不依賴 C 語言的 SQLite 原始程式碼，因此即使 CGO_ENABLED=0（預設在某些環境或 Docker build 中），也能正常編譯與執行。
// 適合用在需要跨平台、無 CGO 的環境，例如 Docker scratch image 或 serverless。

// 為什麼 gorm.io/driver/sqlite 會失敗
// gorm.io/driver/sqlite 預設使用的是 github.com/mattn/go-sqlite3 driver。
// 這個 driver 是 Go 與 SQLite C library 的 binding，需要 CGO 支援。
// 如果你在 CGO_ENABLED=0 的環境下編譯或執行，就會出現錯誤訊息：

import (
	"fmt"
	// https://blog.csdn.net/gitblog_00649/article/details/147110491
	"github.com/glebarez/sqlite"

	// "gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// type OrderModel struct {
// 	DB *gorm.DB
// }
type DBModel struct {
    // 分別是 Order 跟 DB 欄位
	DB *gorm.DB 
	Order OrderModel // *gorm.DB
	
}

// 接收一個 *DBModel 型別指標
// DBModel 通常是一個封裝了多個資料模型的結構體，例如 OrderModel、UserModel 等，負責與資料庫互動。
func InitDB(dataSourceName string)(*DBModel, error){
	
	// 第一個參數 Dialector，用來指定數據庫的類型，像是 mysql / sqlite / postgres 等，db 是由 gorm.Open 回傳的 *gorm.DB 物件， https://zhuanlan.zhihu.com/p/651250516
	db, err := gorm.Open(sqlite.Open(dataSourceName), &gorm.Config{}) // & 表示取址，支持就地修改，如設置日誌或連接池
	
	if err != nil { // 若錯誤發生
		return nil, fmt.Errorf("資料庫連接失敗(ex: 檔案路徑錯誤、目錄不存在或無寫入權限): %w", err)
	}

	err = db.AutoMigrate(&Order{}, &OrderItem{})
	if err != nil {
		return nil, fmt.Errorf("已經成功連接，但執行 schema 遷移時出錯(ex:1. 模型 tag 不支援 SQLite語法 2.更改欄位類型、刪除欄位、check constraint 等，SQLite 限制多，GORM 會重建表格失敗。 3.複合唯一索引、view 依賴、foreign key 順序問題等)")
	}
	
	dbModel := &DBModel{ 
		DB: db, // *gorm.DB 把sqlite的 db gorm物件覆寫
		Order: OrderModel{DB: db}, // 把sqlite的 db gorm物件覆寫在 OrderModel裡的 DB
	}
	return dbModel, nil

}

