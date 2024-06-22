package model

import (
	"fmt"
)

func GetStudent(name string) *User {
	var ret User
	if err := Conn.Table("user").Where("name=?", name).Find(&ret).Error; err != nil {
		fmt.Printf("err:%s", err.Error())
	}
	return &ret
}
func GetAdmin(name string) *Admin {
	var ret Admin
	if err := Conn.Table("admin").Where("name=?", name).Find(&ret).Error; err != nil {
		fmt.Printf("err:%s", err.Error())
	}
	return &ret
}
func CreateUser(user *User) error {
	if err := Conn.Create(user).Error; err != nil {
		fmt.Printf("err:%s", err.Error())
		return err
	}

	return nil
}
