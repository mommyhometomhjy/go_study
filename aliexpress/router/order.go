package router

import (
	"aliexpress/model"
	"aliexpress/vm"
	"net/http"

	"github.com/labstack/echo"
)

func orderIndexHandler(c echo.Context) error {
	vop := vm.OrderViewModelOp{}
	vm := vop.GetVM()

	return c.Render(http.StatusOK, "order/index", &vm)
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