package model

import (
	"time"

	"github.com/jinzhu/gorm"
)

type Order struct {
	gorm.Model
	// 订单号
	OrderNo string

	// 物流方式,物流单号,物流状态,物流花费,包裹重量,签收耗时
	OrderShippingMethod        string
	OrderShippingNo            string
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
