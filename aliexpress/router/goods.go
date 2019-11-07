package router

import (
	"aliexpress/model"
	"strconv"

	"github.com/gin-gonic/gin"
)

//index
func getGoodss(c *gin.Context) {
	goodss := model.GetGoodss()
	c.HTML(200, "goods/index", gin.H{
		"Title":  "产品列表",
		"Goodss": goodss,
	})
}
func goodsNew(c *gin.Context) {
	c.HTML(200, "goods/new", gin.H{
		"Title": "编辑产品",
	})
}

func goodsEdit(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	goods := model.GetGoodsById(id)
	c.HTML(200, "goods/edit", gin.H{
		"Title": "编辑产品",
		"Goods": goods,
	})
}
