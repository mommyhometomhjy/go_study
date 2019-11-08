package model

import (
	"time"
)

type Order struct {
	ID int `gorm:"PRIMARY_KEY"`
	// 订单号
	OrderNo string `gorm:"UNIQUE";form:"OrderNo"`

	// 物流方式,物流单号,物流状态,物流花费,包裹重量,签收耗时
	OrderShippingMethod        string `form:"OrderShippingMethod"`
	OrderShippingNo            string `gorm:"INDEX:shippingnno";form:"OrderShippingNo"`
	OrderShippingStatus        string
	OrderShippingCost          float64
	OrderShippingWeight        float64
	OrderShippingDeliveredDays uint

	// 买家昵称
	OrderBuyer string `form:"OrderBuyer"`

	// 付款时间,付款金额
	OrderPaidTime *time.Time
	OrderMoney    float64 `form:"OrderMoney"`

	// 收件人名称,国家,省份,城市,地址,右边,电话,手机
	OrderReceiverName        string `form:"OrderReceiverName"`
	OrderReceiverCountry     string `form:"OrderReceiverCountry"`
	OrderReceiverProvince    string `form:"OrderReceiverProvince"`
	OrderReceiverCity        string `form:"OrderReceiverCity"`
	OrderReceiverAddress     string `form:"OrderReceiverAddress"`
	OrderReceiverPostCode    string `form:"OrderReceiverPostCode"`
	OrderReceiverTelephone   string `form:"OrderReceiverTelephone"`
	OrderReceiverMobilePhone string `form:"OrderReceiverMobilePhone"`

	// 订单明细
	OrderDetailss []OrderDetails
}

func CountOrderByOrderNo(orderNo string) (total int) {
	var order Order
	db.Where("order_no =?", orderNo).Find(&order).Count(&total)
	return total
}

func GetOrderByShippingNo(no string) []Order {
	var orders []Order
	db.Where("order_shipping_no =? ", no).Find(&orders)
	return orders
}
func UpdateOrder(order *Order) {
	db.Save(order)
}

func UpdateGoodsWeightReferOrder() {
	db.Exec(`update goods
	set goods_weight = (select we from (select avg(weight) as we,goods_no 
	from (select 
		order_shipping_weight as weight,
		goods_no
	from orders 
	join order_details on orders.id = order_details.order_id
	join goods on order_details.goods_id = goods.id
	where order_shipping_weight >0
	group by order_shipping_no 
	having count(order_no)=1)
	group by goods_no) as t2 where t2.goods_no = goods.goods_no limit 1)
	where exists(select * from (select avg(weight) as we,goods_no 
	from (select 
		order_shipping_weight as weight,
		goods_no
	from orders 
	join order_details on orders.id = order_details.order_id
	join goods on order_details.goods_id = goods.id
	where order_shipping_weight >0
	group by order_shipping_no 
	having count(order_no)=1)
	group by goods_no) as t2 where t2.goods_no = goods.goods_no )`)
}

func GetOrders() []Order {
	var orders []Order
	db.Preload("OrderDetailss.Goods").Find(&orders)
	return orders
}

func CreateOrder(order *Order) {
	db.Create(order)
}

func DeleteOrderById(id int) {
	var order Order

	db.First(&order, id)
	db.Delete(&order)
	DeleteOrderDetailsByOrderId(id)
}
func GetOrderById(id int) Order {
	var order Order

	db.Preload("OrderDetailss.Goods").First(&order, id)
	return order
}
