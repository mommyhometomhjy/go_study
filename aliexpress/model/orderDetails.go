package model

import (
	"github.com/jinzhu/gorm"
)

type OrderDetails struct {
	gorm.Model
	// 关联的订单id
	Order
	OrderId uint

	Goods
	GoodsId uint
	//当前产品的数量
	Number uint
}

func CreateOrderDetails(orderDetails *OrderDetails) {
	db.Create(orderDetails)
}
