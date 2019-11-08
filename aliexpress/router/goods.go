package router

import (
	"aliexpress/model"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/gin-gonic/gin"
)

//index
func getGoodss(c *gin.Context) {
	sp := c.DefaultQuery("page", "1")
	p, _ := strconv.Atoi(sp)
	no := c.PostForm("no")
	if no == "" {
		no = c.Query("no")
	}
	if no != "" {
		no = strings.ToUpper(no)
	}
	goodss, total := model.SearchGoodsByGoodsNo(no, p, 10)
	page := BasePageViewModel{}
	page.SetBasePageViewModel(total, p, 10)
	page.SetPrevAndNextPage()
	c.HTML(200, "goods/index", gin.H{
		"Title":  "产品列表",
		"Goodss": goodss,
		"Page":   page,
		"Search": no,
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
	c.Redirect(http.StatusMovedPermanently, "/goods/index?no="+goodsno)
}

func parseStandardShippingCost(c *gin.Context) {
	model.DeleteStandShippingCost()

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
		return
	}
	rows := f.GetRows("Worksheet1")
	if err != nil {
		fmt.Println(err)
		return
	}
	for index, row := range rows {
		if index == 0 {
			continue
		}
		weight := row[0]
		price, _ := strconv.ParseFloat(row[1], 64)
		s := model.StandShippingCost{
			Weight: weight,
			Price:  price,
		}
		model.CreateShippingCost(&s)
	}
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

func exportsellpricechanged(c *gin.Context) {
	goodss := model.GetGoodsSellPriceChanged()

	c.Writer.Header().Add("Content-type", "application/octet-stream")
	c.Writer.Header().Add("content-disposition", "attachment; filename=\"价格变动.xlsx\"")
	f := excelize.NewFile()

	f.SetCellValue("Sheet1", "A1", "货号")
	f.SetCellValue("Sheet1", "B1", "成本价")
	f.SetCellValue("Sheet1", "C1", "重量")
	f.SetCellValue("Sheet1", "D1", "售价")
	f.SetCellValue("Sheet1", "E1", "上次售价")

	for index, goods := range goodss {
		sheeti := strconv.Itoa(index + 2)
		f.SetCellValue("Sheet1", "A"+sheeti, goods.GoodsNo)
		f.SetCellValue("Sheet1", "B"+sheeti, goods.GoodsPrice)
		f.SetCellValue("Sheet1", "C"+sheeti, goods.GoodsWeight)
		f.SetCellValue("Sheet1", "D"+sheeti, goods.GoodsSellPrice)
		f.SetCellValue("Sheet1", "E"+sheeti, goods.GoodsLastSellPrice)
		model.UpdateGoods(&goods)
	}
	f.Write(c.Writer)

}
