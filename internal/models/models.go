package models

import (
	"fmt"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type DBModel struct {
	Order OrderModel
	DB *gorm.DB
}

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
		DB: db,
		Order: OrderModel{DB: db},
	}
	return dbModel, nil

}

