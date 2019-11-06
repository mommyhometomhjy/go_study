package router

import (
	"aliexpress/model"
	"aliexpress/vm"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func orderIndexHandler(c *gin.Context) {
	vop := vm.OrderViewModelOp{}
	vm := vop.OrderGetIndexVM()
	c.HTML(http.StatusOK, "order/index", &vm)

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

	vop := vm.OrderViewModelOp{}
	vm := vop.OrderGetNewVM()
	c.HTML(http.StatusOK, "order/new", &vm)
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

	vop := vm.OrderViewModelOp{}
	vm := vop.OrderGetIndexVM()
	c.HTML(http.StatusOK, "order/index", &vm)

}

func orderDelete(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	model.DeleteOrderById(id)
	c.String(200, "successed")
}

func orderEditHandler(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	vop := vm.OrderViewModelOp{}
	vm := vop.OrderGetEditVM(id)
	c.HTML(http.StatusOK, "order/edit", &vm)

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
