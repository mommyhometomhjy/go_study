package model

import (
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/jinzhu/gorm"
)

type Order struct {
	gorm.Model
	// 订单号
	OrderNo string `gorm:"UNIQUE"`

	// 物流方式,物流单号,物流状态,物流花费,包裹重量,签收耗时
	OrderShippingMethod        string
	OrderShippingNo            string `gorm:"INDEX:shippingnno"`
	OrderShippingStatus        string
	OrderShippingCost          float64
	OrderShippingWeight        float64
	OrderShippingDeliveredDays uint

	// 买家昵称
	OrderBuyer string

	// 付款时间,付款金额
	OrderPaidTime *time.Time
	OrderMoney    float64

	// 收件人名称,国家,省份,城市,地址,右边,电话,手机
	OrderReceiverName        string
	OrderReceiverCountry     string
	OrderReceiverProvince    string
	OrderReceiverCity        string
	OrderReceiverAddress     string
	OrderReceiverPostCode    string
	OrderReceiverTelephone   string
	OrderReceiverMobilePhone string

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
	db.Preload("OrderDetailss").Find(&orders)
	return orders
}
func ParseOrderExcel(r io.Reader) (err error) {
	f, err := excelize.OpenReader(r)
	if err != nil {
		fmt.Println(err)
		return
	}

	rows := f.GetRows("sheet1")
	if err != nil {
		fmt.Println(err)
		return
	}
	for index, row := range rows {
		//跳过第一行
		if index == 0 {
			continue
		}
		//跳过没有单号的行
		if row[24] == "" {
			continue
		}
		//跳过没有sku的行
		if row[11] == "" {
			continue
		}
		//跳过已经更新了的订单
		if CountOrderByOrderNo(row[0]) > 0 {
			continue
		}

		//跳过未付款的订单
		if row[5] == "" {
			continue
		}
		// fmt.Println(row[5])
		//获取时间
		paidTime, err := time.Parse("2006-01-02 15:04:05", row[5]+":00")
		if err != nil {
			fmt.Println(err)
			return err
		}
		//获取订单金额
		ordermonkey, _ := strconv.ParseFloat(strings.Replace(row[8], "US $", "", -1), 64)

		skuNumArray := strings.Split(row[11], ";")
		var orderDetailss []OrderDetails

		for _, skuNum := range skuNumArray {
			goods := FindGoodsByGoodsNo(strings.Split(skuNum, " * ")[0])

			number, _ := strconv.Atoi(strings.Split(skuNum, " * ")[1])
			orderDetails := OrderDetails{
				Goods:  goods,
				Number: uint(number),
			}
			orderDetailss = append(orderDetailss, orderDetails)
		}

		order := Order{
			//订单号
			OrderNo: row[0],

			// 物流方式,物流单号,物流状态,物流花费,包裹重量,签收耗时
			OrderShippingMethod:        strings.Trim(strings.Split(row[24], ":")[0], " \n"),
			OrderShippingNo:            strings.Trim(strings.Split(row[24], ":")[1], " \n"),
			OrderShippingStatus:        row[1],
			OrderShippingCost:          0.0,
			OrderShippingWeight:        0.0,
			OrderShippingDeliveredDays: 0,

			// 买家昵称
			OrderBuyer: row[3],

			// 付款时间,付款金额
			OrderPaidTime: &paidTime,
			OrderMoney:    ordermonkey,

			// 收件人名称,国家,省份,城市,地址,右边,电话,手机
			OrderReceiverName:        row[14],
			OrderReceiverCountry:     row[15],
			OrderReceiverProvince:    row[16],
			OrderReceiverCity:        row[17],
			OrderReceiverAddress:     row[18],
			OrderReceiverPostCode:    row[19],
			OrderReceiverTelephone:   row[20],
			OrderReceiverMobilePhone: row[21],
			OrderDetailss:            orderDetailss,
		}
		db.Create(&order)
		// fmt.Println(order.OrderNo, len(order.OrderDetailss), order.OrderDetailss[0].Goods.GoodsNo)

	}
	return nil
}
func ParseShippingCost() {

	f, err := excelize.OpenFile("cmd/shipfundList.xlsx")
	if err != nil {
		fmt.Println(err)
		return
	}

	rows := f.GetRows("sheet1")
	if err != nil {
		fmt.Println(err)
		return
	}
	for index, row := range rows {
		if index == 0 {
			continue
		}
		orders := GetOrderByShippingNo(row[1])
		// fmt.Println(len(orders), row[1])
		cost, _ := strconv.ParseFloat(row[7], 64)
		weight, _ := strconv.ParseFloat(row[5], 64)
		if weight < 10 {
			weight *= 1000
		}
		for _, order := range orders {
			order.OrderShippingCost = cost
			order.OrderShippingWeight = weight
			UpdateOrder(&order)

		}
	}

	UpdateGoodsWeightReferOrder()
	UpdateGoodsPrice()
}
