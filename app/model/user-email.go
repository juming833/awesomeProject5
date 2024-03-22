package model

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
)

// GetUserByEmail 通过邮箱获取用户信息
func GetUserByEmail(email string) *User {
	var ret User
	if err := Conn.Table("user").Where("email = ?", email).Find(&ret).Error; err != nil {
		fmt.Printf("err:%s", err.Error())
	}
	return &ret
}

func GetUserByPhone(phone string) (int, error) {
	var userId int
	err := Conn.Transaction(func(db *gorm.DB) error {
		err := db.Table("user").Where("phone = ?", phone).Select("id").Scan(&userId).Error
		if err != nil {
			return err
		}
		return nil
	})
	// 查询用户信息
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, errors.New("用户不存在")
		}
		return 0, err
	}

	return userId, nil
}
