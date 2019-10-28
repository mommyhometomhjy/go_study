package model

import (
	"github.com/jinzhu/gorm"
)

type OrderDetails struct {
	gorm.Model
	// 关联的订单id
	OrderId uint

	// 关联的产品id
	GoodsId uint

	//当前产品的数量
	Number uint
}
