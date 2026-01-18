package models

/*
來源: https://github.com/glebarez/sqlite?tab=readme-ov-file#how-is-this-better-than-standard-gorm-sqlite-driver
基於 cgo：官方的 gorm.io/driver/sqlite 是透過 Go 與 SQLite C 原始碼的綁定（cgo）。
限制：
需要安裝 C 編譯器才能編譯與執行程式。
SQLite 的某些功能（例如 JSON 支援）必須在編譯時啟用，因此每次執行 go run、go test 等指令時都要加上正確的 build tags。
因為需要 C 編譯器，無法在精簡的容器（例如 golang-alpine）中建置。
在 GCP（Google Cloud Platform）上無法建置，因為 GCP 不允許執行 gcc。

glebarez/sqlite 的優勢
純 Go 實作：這個 driver 基於 cznic/sqlite，它是將 SQLite C 原始碼 AST 轉換成 Go 程式碼。

好處：
不需要 C 編譯器，跨平台部署更方便。
可以在精簡容器（如 golang-alpine）或 GCP 上使用。
本質上仍是原始 SQLite 的邏輯，只是用 Go 語言重寫。
*/

import (
	"fmt"
	// sqlite3 -header -column data/orders.db "SELECT * FROM orders;"
	// https://blog.csdn.net/gitblog_00649/article/details/147110491
	"github.com/glebarez/sqlite"
	// 不論是 gorm.io/driver/sqlite 還是 github.com/glebarez/sqlite，生成的 .db 檔案都能用 sqlite3 -header -column data/orders.db "SELECT * FROM orders;" 查詢。差別只在 Go driver 的實作方式，跟 CLI 指令無關。
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

