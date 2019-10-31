package router

import (
	"aliexpress/vm"
	"net/http"

	"github.com/labstack/echo"
)

func indexHandler(c echo.Context) error {
	vop := vm.IndexViewModelOp{}
	vm := vop.GetVM()

	return c.Render(http.StatusOK, "index", &vm)
}
