package cmd

import (
	"aliexpress/model"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/360EntSecGroup-Skylar/excelize"
)

func ParseExcel() {
	f, err := excelize.OpenFile("cmd/导出订单.xlsx")
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

		//跳过已经更新了的订单
		if model.CountOrderByOrderNo(row[0]) > 0 {
			continue
		}
		// fmt.Println(row[5])
		//获取时间
		paidTime, err := time.Parse("2006-01-02 15:04:05", row[5]+":00")
		if err != nil {
			fmt.Println(err)
			return
		}
		//获取订单金额
		ordermonkey, _ := strconv.ParseFloat(strings.Replace(row[8], "US $", "", -1), 64)

		order := model.Order{
			//订单号
			OrderNo: row[0],

			// 物流方式,物流单号,物流状态,物流花费,包裹重量,签收耗时
			OrderShippingMethod:        strings.Split(row[24], ":")[0],
			OrderShippingNo:            strings.Split(row[24], ":")[1],
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
		}
		model.CreateOrder(&order)

		skuNumArray := strings.Split(row[11], ";")
		for _, skuNum := range skuNumArray {
			goods := model.FindGoodsByGoodsNo(strings.Split(skuNum, " * ")[0])

			number, _ := strconv.Atoi(strings.Split(skuNum, " * ")[1])
			orderDetails := model.OrderDetails{
				OrderId: order.ID,
				GoodsId: goods.ID,
				Number:  uint(number),
			}
			model.CreateOrderDetails(&orderDetails)
		}
	}

}
