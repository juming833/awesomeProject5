package models

import "gorm.io/gorm"

type UserModel struct {
	gorm.Model
	Username     string `gorm:"size:32" json:"username"`
	Password     string `gorm:"size:64" json:"password"`
	CaptchaId    string `json:"captchaId"`
	CaptchaValue string `json:"captchaValue"`
}
