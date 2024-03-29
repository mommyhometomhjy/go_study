package model

type OrderDetails struct {
	ID int `gorm:"PRIMARY_KEY"`
	// 关联的订单id
	OrderId uint

	Goods   Goods
	GoodsId uint
	//当前产品的数量
	Number uint `form:"Number"`
}

func CreateOrderDetails(orderDetails *OrderDetails) {
	db.Create(orderDetails)
}
func DeleteOrderDetailsByOrderId(orderId int) {

	db.Where("order_id=? ", orderId).Delete(OrderDetails{})

}

func OrderDetailsHasGoods(id int) bool {
	var o []OrderDetails
	var count int
	db.Where("goods_id = ?", id).Find(&o).Count(&count)
	if count > 0 {
		return true
	}
	return false
}
