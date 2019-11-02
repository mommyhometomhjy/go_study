package router

import (
	"aliexpress/model"
	"aliexpress/vm"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo"
)

func orderIndexHandler(c echo.Context) error {
	vop := vm.OrderViewModelOp{}
	vm := vop.OrderGetIndexVM()
	err := c.Render(http.StatusOK, "order/index", &vm)
	if err != nil {
		fmt.Println(err)
	}
	return err
}

func orderImportExcel(c echo.Context) error {
	file, err := c.FormFile("file")
	if err != nil {
		return err
	}
	// fmt.Println(file.Filename)
	src, err := file.Open()
	if err != nil {
		// fmt.Println(err)
		return err
	}
	defer src.Close()
	model.ParseOrderExcel(src)
	return orderIndexHandler(c)

}

func orderNewHandler(c echo.Context) error {

	vop := vm.OrderViewModelOp{}
	vm := vop.OrderGetNewVM()
	return c.Render(http.StatusOK, "order/new", &vm)
}

func orderCreate(c echo.Context) error {
	type GoodsNoAndNumber struct {
		GoodsNO []string `form:"GoodsNo"`
		Number  []uint   `form:"Number"`
	}
	var order model.Order
	var goodss GoodsNoAndNumber
	var orderDetailss []model.OrderDetails

	if err := c.Bind(&goodss); err != nil {
		return err
	}
	for index, goodsNo := range goodss.GoodsNO {
		goods := model.FindGoodsByGoodsNo(goodsNo)

		orderDetailss = append(orderDetailss, model.OrderDetails{
			Goods:  goods,
			Number: goodss.Number[index],
		})
	}

	if err := c.Bind(&order); err != nil {
		return err
	}

	order.OrderDetailss = orderDetailss
	t := time.Now()
	order.OrderPaidTime = &t
	model.CreateOrder(&order)

	vop := vm.OrderViewModelOp{}
	vm := vop.OrderGetIndexVM()
	err := c.Render(http.StatusOK, "order/index", &vm)
	return err
}

func orderDelete(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))

	model.DeleteOrderById(id)
	return c.String(200, "successed")
}

func orderEditHandler(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	vop := vm.OrderViewModelOp{}
	vm := vop.OrderGetEditVM(id)
	err := c.Render(http.StatusOK, "order/edit", &vm)
	if err != nil {
		fmt.Println(err)
	}
	return err

}

func orderUpdate(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	order := model.GetOrderById(id)

	if err := c.Bind(&order); err != nil {
		return err
	}
	type GoodsNoAndNumber struct {
		GoodsNO []string `form:"GoodsNo"`
		Number  []uint   `form:"Number"`
	}

	var goodss GoodsNoAndNumber
	var orderDetailss []model.OrderDetails

	if err := c.Bind(&goodss); err != nil {
		return err
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
	return orderEditHandler(c)
}
