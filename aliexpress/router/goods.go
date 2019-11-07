package router

import (
	"aliexpress/model"
	"fmt"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

//index
func getGoodss(c *gin.Context) {
	p := 1
	sp := c.Query("page")
	if sp != "" {
		p, _ = strconv.Atoi(sp)
	}
	goodss, page := model.GetGoodss(p)
	c.HTML(200, "goods/index", gin.H{
		"Title":  "产品列表",
		"Goodss": goodss,
		"Page":   page,
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

func goodsCreate(c *gin.Context) {
	goodsno := c.PostForm("goodsno")
	styles := c.PostForm("styles")
	sizes := c.PostForm("sizes")
	weight, _ := strconv.ParseFloat(c.PostForm("weight"), 64)
	price, _ := strconv.ParseFloat(c.PostForm("price"), 64)
	link := c.PostForm("link")

	switch {
	case styles != "" && sizes != "":

		var styleArr, sizeArr []string
		styleArr = strings.Split(styles, ",")
		sizeArr = strings.Split(sizes, ",")
		for _, style := range styleArr {
			for _, size := range sizeArr {
				goods := model.Goods{
					GoodsNo:           goodsno + "-" + style + "-" + size,
					GoodsPrice:        price,
					GoodsWeight:       weight,
					GoodsSupplierLink: link,
				}
				fmt.Println(goods)
				model.CreateGoods(&goods)
			}
		}
	case styles == "" && sizes == "":
		goods := model.Goods{
			GoodsNo:           goodsno,
			GoodsPrice:        price,
			GoodsWeight:       weight,
			GoodsSupplierLink: link,
		}
		model.CreateGoods(&goods)
	default:
		variant := styles + sizes
		vArr := strings.Split(variant, ",")
		for _, vari := range vArr {
			goods := model.Goods{
				GoodsNo:           goodsno + "-" + vari,
				GoodsPrice:        price,
				GoodsWeight:       weight,
				GoodsSupplierLink: link,
			}
			model.CreateGoods(&goods)
		}
	}
	model.UpdateGoodsPrice()
	getGoodss(c)
}

func parseStandardShippingCost(c *gin.Context) {

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

	model.ParseStandardShippingExcel(src)
	model.UpdateGoodsPrice()
	getGoodss(c)
}

func goodsDelete(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	count := model.OrderDetailsHasGoods(id)
	if count {
		c.String(200, "存在订单,禁止删除")
	} else {
		model.DeleteGoodsById(id)
		c.String(200, "删除成功")
	}

}

func goodsUpdate(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	weight, _ := strconv.ParseFloat(c.PostForm("weight"), 64)
	price, _ := strconv.ParseFloat(c.PostForm("price"), 64)
	link := c.PostForm("link")
	goods := model.GetGoodsById(id)
	goods.GoodsPrice = price
	goods.GoodsWeight = weight
	goods.GoodsSupplierLink = link
	model.UpdateGoods(&goods)
	goodsEdit(c)
}
