package router

import (
	"aliexpress/model"
	"aliexpress/vm"
	"fmt"
	"net/http"

	"github.com/labstack/echo"
)

func orderIndexHandler(c echo.Context) error {
	vop := vm.OrderViewModelOp{}
	vm := vop.GetVM()
	err := c.Render(http.StatusOK, "order/index", &vm)
	// fmt.Println(err)
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
	vm := vop.GetVM()
	vm.SetTitle("新建订单")
	return c.Render(http.StatusOK, "order/new", &vm)
}

func orderCreate(c echo.Context) error {
	var order model.Order
	goodsno := c.FormValue("GoodsNo")
	fmt.Println(goodsno)
	if err := c.Bind(&order); err != nil {
		return err
	}
	model.CreateOrder(&order)
	vop := vm.OrderViewModelOp{}
	vm := vop.GetVM()
	err := c.Render(http.StatusOK, "order/index", &vm)
	return err
}
