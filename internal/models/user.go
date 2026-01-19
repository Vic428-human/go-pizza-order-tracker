package models

import "gorm.io/gorm"

type User struct { // 定義一個 User struct，對應到資料庫中的 users table。
	ID       string `gorm:"primaryKey"`
	Username string `gorm:"uniqueIndex;not null"`
	Password string `gorm:"not null"`
}

type UserModel struct {
	DB *gorm.DB // 持有 *gorm.DB，用來執行資料庫操作。
}

// 這是一個方法，綁定在 UserModel 上。
// user, password => 參數
// *User：如果驗證成功，回傳該使用者資料。 error：如果失敗，回傳錯誤。
func (u *UserModel) AuthenticateUser(username, password string) (*User, error) {
	var user User
  
	// 用 GORM 查詢 users table，找出符合 username 的第一筆資料。
	if err:= u.DB.Where("username=?", username).First(&user).Error; err != {
		// 如果沒有找到紀錄，GORM 會回傳這個錯誤。你可以選擇回傳 nil 代表不存在。
		if err == gorm.ErrRecordNotFound{
			return nil,nil;
		}
		return nil, err;
	}
	

	// 密碼不應該存明文，通常會用 bcrypt 或其他 hash 函式加密。

	
	return &User,nil // 成功找到的話，就回傳 *User。

}
