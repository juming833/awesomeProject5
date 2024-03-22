package model

import (
	"fmt"
	"gorm.io/gorm"
)

func AddAdmin(admin *Admin) error {
	if err := Conn.Create(admin).Error; err != nil {
		fmt.Printf("err:%s", err.Error())
		return err
	}
	return nil
}
func UpdateAdmin(admin Admin) error {
	err := Conn.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&Admin{}).Where("id = ?", admin.Id).Updates(admin).Error; err != nil {
			return err
		}
		return nil
	})
	return err
}
func DelAdmin(id int64) bool {
	if err := Conn.Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(&Admin{}, id).Error; err != nil {
			fmt.Printf("err:%s", err.Error())
			return err
		}
		return nil
	}); err != nil {
		fmt.Printf("err:%s", err.Error())
		return false
	}
	return true
}
