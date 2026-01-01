package models

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type DBModel struct {
	DB *gorm.DB
}

func InitDB(dataSourceName string)(*DBModel, error){
	
	// 創建數據庫連接
	// https://zhuanlan.zhihu.com/p/651250516
	// 第一個參數 Dialector，用來指定數據庫的類型，像是 mysql / sqlite / postgres 等等
	// db 是由 gorm.Open 回傳的 *gorm.DB 物件
	db, err := gorm.Open(sqlite.Open(dataSourceName), &gorm.Config{}) // & 表示取址，支持就地修改，如設置日誌或連接池
	
	if err != nil {
		return nil, err
	}

	// 將 gorm 的實例設定到 DBModel，讓該結構體能直接暴露資料庫的方法 => dbModel.DB.Create(&user)
	dbModel := &DBModel{ 
		DB: db,
	}


	return dbModel, nil
}

