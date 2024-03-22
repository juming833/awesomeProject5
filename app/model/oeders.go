package model

import (
	"errors"
	"gorm.io/gorm"
	"time"
)

func IsOrderOpen(orderNo string) bool {
	// 检查订单状态，返回订单是否处于打开状态
	var order Orders
	err := Conn.Transaction(func(db *gorm.DB) error {
		if err := db.Table("orders").Where("order_no = ?", orderNo).First(&order).Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		// 处理查询错误
		return false
	}
	return order.Status == "OPEN"
}

func UpdateOrderStatus(orderNo string, status string) error {
	// 更新订单状态为指定状态
	err := Conn.Transaction(func(db *gorm.DB) error {
		if err := db.Table("orders").Where("order_no = ?", orderNo).Update("status", status).Error; err != nil {
			return err
		}
		return nil
	})
	return err
}
func CreateOrder(orderNo string) error {
	var order Orders
	err := Conn.Transaction(func(db *gorm.DB) error {
		if err := db.Table("orders").Where("order_no = ?", orderNo).First(&order).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				// 订单号为orderNo的记录不存在，创建新记录
				order := Orders{
					OrderNo:     orderNo,
					Status:      "OPEN",
					CreatedTime: time.Now(),
					UpdatedTime: time.Now(),
				}
				if err := db.Create(&order).Error; err != nil {
					return err
				}
				return nil
			}
			return err
		}
		return nil
	})
	return err
}
