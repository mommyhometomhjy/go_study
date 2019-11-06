package router

import (
	"aliexpress/vm"
	"net/http"

	"github.com/gin-gonic/gin"
)

func indexHandler(c *gin.Context) {
	vop := vm.IndexViewModelOp{}
	vm := vop.GetVM()
	c.HTML(http.StatusOK, "index", &vm)
}
