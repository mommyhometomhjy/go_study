package router

import (
	"aliexpress/model"
	"fmt"
	"strconv"
	"time"

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
	model.ParseOrderExcel(src)
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
		goods := model.FindGoodsByGoodsNo(goodsNo)

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
		goods := model.FindGoodsByGoodsNo(goodsNo)

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
