package models

import (
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct { // 定義一個 User struct，對應到資料庫中的 users table。
	ID       string `gorm:"primaryKey"`
	Username string `gorm:"uniqueIndex;not null"`
	Password string `gorm:"not null"`
}

type UserModel struct {
	DB *gorm.DB // 持有 *gorm.DB，用來執行資料庫操作。
}

func HashPassword(password string) (string, error) {
	// 以给定的Cost返回密码的bcrypt哈希。如果给定的成本小于MinCost，则将成本设置为DefaultCost（10）
    bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
    return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	// 用于比对bcrypt哈希字符串和提供的密码明文文本是否匹配
    if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)); err != nil{
		return nil, errors.New("Invalid Credientials")
	}
	return err == nil
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
	

	// 密碼不應該存明文，通常會用 bcrypt 或其他 hash 函式加密 => https://kryiea.github.io/back-end/go/framework/encipherframework/bcrypt%E5%8A%A0%E5%AF%86.html#_4-1-bcrypt%E5%8C%85%E4%BB%8B%E7%BB%8D
	hash, _ := HashPassword(password) 
	fmt.Fprintln(w, "Hash:    ", hash) // Hash:     $2a$14$Ael8nW7UF/En/iI7LGdyBuaIO8VREbL2CAShRN0EUQHqtmOHXh.XK
	match := CheckPasswordHash(password, hash) 
	fmt.Fprintln(w, "Match:   ", match) // Match:    true
	
	return &User, nil // 成功找到的話，就回傳 *User。

}
