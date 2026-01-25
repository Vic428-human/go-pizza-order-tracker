package models

import (
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct { // 定義一個 User struct，對應到資料庫中的 users table。
	ID       uint   `gorm:"primaryKey"`
	Username string `gorm:"uniqueIndex;not null"`
	Password string `gorm:"not null"`
}

type UserModel struct {
	DB *gorm.DB // 持有 *gorm.DB，用來執行資料庫操作。
}

func GenerateHashPassword(password string) (string, error) {
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashPassword), nil
}

func CompareHashAndPassword(hashPassword string, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashPassword), []byte(password))
	return err == nil
}

// 這方法綁定在 UserModel 上
// user, password => 參數
// *User：如果驗證成功，回傳該使用者資料。 error：如果失敗，回傳錯誤。
func (u *UserModel) AuthenticateUser(username, password string) (*User, error) {
	var user User

	// 查詢符合 username 的第一筆資料
	if err := u.DB.Where("username = ?", username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("使用者不存在")
		}

		return nil, fmt.Errorf("系統錯誤，請稍後再試")
	}
	fmt.Printf("====>從資料庫取得hash後的密碼:%s\n", user.Password)
	if CompareHashAndPassword(user.Password, password) {
		fmt.Println("===>匹配")
		return &user, nil
	} else {
		fmt.Println("===>密碼匹配錯誤")
		return nil, fmt.Errorf("密碼匹配錯誤")
	}
}

func (u *UserModel) GetUserByID(id string) (*User, error) {
	var user User
	if err := u.DB.First(&user, "id = ?", id).Error; err != nil { // SELECT * FROM users WHERE id = "1b74413f-f3b8-409f-ac47-e8c062e3472a";
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // 資料庫裡已經沒有這個使用者
		}
		return nil, err // 資料庫錯誤
	}
	return &user, nil
}
