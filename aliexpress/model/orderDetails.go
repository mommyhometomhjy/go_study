package model

import (
	"github.com/jinzhu/gorm"
)

type OrderDetails struct {
	gorm.Model
	// 关联的订单id
	Order Order `gorm:"association_foreignkey:ID"`

	Goods Goods `gorm:"association_foreignkey:ID"`
	//当前产品的数量
	Number uint
}

func CreateOrderDetails(orderDetails *OrderDetails) {
	db.Create(orderDetails)
}
