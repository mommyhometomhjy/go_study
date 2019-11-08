package router

import (
	"aliexpress/model"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/gin-gonic/gin"
)

func orderIndexHandler(c *gin.Context) {

	orders := model.GetOrders()
	c.HTML(200, "order/index", gin.H{
		"Title":  "订单列表",
		"Orders": orders,
	})

}

func orderImportExcel(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		fmt.Println(err)
	}
	// fmt.Println(file.Filename)
	src, err := file.Open()
	if err != nil {
		// fmt.Println(err)
		fmt.Println(err)
	}
	defer src.Close()

	f, err := excelize.OpenReader(src)
	if err != nil {
		fmt.Println(err)

	}

	rows := f.GetRows("sheet1")
	if err != nil {
		fmt.Println(err)

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
		if model.CountOrderByOrderNo(row[0]) > 0 {
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
		}
		//获取订单金额
		ordermonkey, _ := strconv.ParseFloat(strings.Replace(row[8], "US $", "", -1), 64)

		skuNumArray := strings.Split(row[11], ";")
		var orderDetailss []model.OrderDetails

		for _, skuNum := range skuNumArray {
			goods := model.GetGoodsByGoodsNo(strings.Split(skuNum, " * ")[0])

			number, _ := strconv.Atoi(strings.Split(skuNum, " * ")[1])
			orderDetails := model.OrderDetails{
				Goods:  goods,
				Number: uint(number),
			}
			orderDetailss = append(orderDetailss, orderDetails)
		}

		order := model.Order{
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
		model.CreateOrder(&order)
		// fmt.Println(order.OrderNo, len(order.OrderDetailss), order.OrderDetailss[0].Goods.GoodsNo)

	}
	orderIndexHandler(c)

}

func orderNewHandler(c *gin.Context) {

	c.HTML(200, "order/new", gin.H{
		"Title": "新建订单",
	})
}

func orderCreate(c *gin.Context) {
	type GoodsNoAndNumber struct {
		GoodsNO []string `form:"GoodsNo"`
		Number  []uint   `form:"Number"`
	}
	var order model.Order
	var goodss GoodsNoAndNumber
	var orderDetailss []model.OrderDetails

	if err := c.Bind(&goodss); err != nil {
		fmt.Println(err)
	}
	for index, goodsNo := range goodss.GoodsNO {
		goods := model.GetGoodsByGoodsNo(goodsNo)

		orderDetailss = append(orderDetailss, model.OrderDetails{
			Goods:  goods,
			Number: goodss.Number[index],
		})
	}

	if err := c.Bind(&order); err != nil {
		fmt.Println(err)
	}

	order.OrderDetailss = orderDetailss
	t := time.Now()
	order.OrderPaidTime = &t
	model.CreateOrder(&order)

	orderIndexHandler(c)

}

func orderDelete(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	model.DeleteOrderById(id)
	c.String(200, "successed")
}

func orderEditHandler(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	order := model.GetOrderById(id)

	c.HTML(200, "order/edit", gin.H{
		"Titile": "编辑订单",
		"Order":  order,
	})

}

func orderUpdate(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	order := model.GetOrderById(id)

	if err := c.Bind(&order); err != nil {
		fmt.Println(err)
	}
	type GoodsNoAndNumber struct {
		GoodsNO []string `form:"GoodsNo"`
		Number  []uint   `form:"Number"`
	}

	var goodss GoodsNoAndNumber
	var orderDetailss []model.OrderDetails

	if err := c.Bind(&goodss); err != nil {
		fmt.Println(err)
	}
	for index, goodsNo := range goodss.GoodsNO {
		goods := model.GetGoodsByGoodsNo(goodsNo)

		orderDetailss = append(orderDetailss, model.OrderDetails{
			Goods:  goods,
			Number: goodss.Number[index],
		})
	}
	model.DeleteOrderDetailsByOrderId(id)
	order.OrderDetailss = orderDetailss
	model.UpdateOrder(&order)
	orderEditHandler(c)
}
