package models

import (
	"errors"

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
	// 用于比对 bcrypt 哈希字符串和提供的密码明文文本是否匹配
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// 這是一個方法，綁定在 UserModel 上。
// user, password => 參數
// *User：如果驗證成功，回傳該使用者資料。 error：如果失敗，回傳錯誤。
func (u *UserModel) AuthenticateUser(username, password string) (*User, error) {
	var user User

	// 查詢符合 username 的第一筆資料
	if err := u.DB.Where("username = ?", username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // 找不到就回傳 nil，不算錯誤
		}
		return nil, err
	}

	// 驗證密碼：比對使用者輸入的 password 與資料庫中的 hash
	match := CheckPasswordHash(password, user.Password)
	if match {
		return &user, nil
	}

	// 如果密碼不正確，回傳 nil 與錯誤
	return nil, errors.New("invalid credentials")
}

func (u *UserModel) GetUserByID(id string) (*User, error) {
	var user User
	if err := u.DB.First(&user, "id = ?", id).Error; err != nil { // SELECT * FROM users WHERE id = "1b74413f-f3b8-409f-ac47-e8c062e3472a";
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // 找不到就回傳 nil，不算錯誤
		}
		return nil, err
	}
	return &user, nil
}
