package model

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/360EntSecGroup-Skylar/excelize"
)

type Goods struct {
	ID int `gorm:"PRIMARY_KEY"`
	// 产品sku
	GoodsNo string `gorm:"UNIQUE";form:"GoodsNo"`

	// 速卖通id
	AliexpressId string

	// 产品进价
	GoodsPrice float64
	// 单个包裹重量
	GoodsWeight float64
	//库存
	GoodsStock uint
	//售价
	GoodsSellPrice float64

	//上次售价
	GoodsLastSellPrice float64
}

func FindGoodsByGoodsNo(goodsNo string) Goods {
	var goods Goods
	db.FirstOrCreate(&goods, Goods{GoodsNo: strings.ToUpper(goodsNo)})
	return goods
}
func UpdateGoods(goods *Goods) {
	db.Save(goods)
}
func UpdateGoodsPrice() {
	var goodss []Goods
	percent, exchange := 0.857, 7.0
	db.Where("goods_weight >0").Find(&goodss)
	for _, goods := range goodss {

		w := fmt.Sprintf("%d", int(math.Ceil(goods.GoodsWeight/10.0)*10))

		standShippingCost := GetPriceByWeight(w)
		goods.GoodsLastSellPrice = goods.GoodsSellPrice
		goods.GoodsSellPrice = math.Ceil((goods.GoodsPrice+3+standShippingCost)/percent/exchange) - 0.01
		db.Save(goods)
	}
}

func GetGoodsIncludeSellPriceAndAliexpressId() []Goods {
	var goodss []Goods
	db.Where("goods_sell_price <>goods_last_sell_price and aliexpress_id <>''").Order("goods_no").Find(&goodss)
	return goodss
}
func ParseGoodsExcel() {
	f, err := excelize.OpenFile("cmd/店小秘导出商品.xlsx")
	if err != nil {
		fmt.Println(err)
		return
	}

	rows := f.GetRows("产品列表")
	if err != nil {
		fmt.Println(err)
		return
	}
	for index, row := range rows {
		if index == 0 {
			continue
		}
		goods := FindGoodsByGoodsNo(row[1])
		goods.AliexpressId = row[0]
		stock, _ := strconv.Atoi(row[3])
		goods.GoodsStock = uint(stock)
		UpdateGoods(&goods)
	}
}
func ParseStandShippingCost() {
	f, err := excelize.OpenFile("cmd/标准运费.xlsx")
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

		price, err := strconv.ParseFloat(row[1], 64)
		if err != nil {
			fmt.Println(err)
			return
		}
		shippingCost := GetShippingCostByWeight(row[0])
		shippingCost.Price = price
		UpdateShippingCost(&shippingCost)
	}
}
func ExportGoodsIncludePrice() {
	goodss := GetGoodsIncludeSellPriceAndAliexpressId()
	f := excelize.NewFile()
	f.SetCellValue("Sheet1", "A1", `*产品ID
	（必填）
	`)
	f.SetCellValue("Sheet1", "B1", `*商品编码
	（必填）`)
	f.SetCellValue("Sheet1", "C1", `*价格
	（必填）
	`)
	f.SetCellValue("Sheet1", "D1", `*库存
	（必填)
	`)

	for index, goods := range goodss {
		f.SetCellValue("Sheet1", fmt.Sprint("A%d", index+2), goods.AliexpressId)
		f.SetCellValue("Sheet1", fmt.Sprint("B%d", index+2), goods.GoodsNo)
		f.SetCellValue("Sheet1", fmt.Sprint("C%d", index+2), goods.GoodsSellPrice)
		f.SetCellValue("Sheet1", fmt.Sprint("D%d", index+2), goods.GoodsStock)
	}

	err := f.SaveAs("./Book1.xlsx")
	if err != nil {
		fmt.Println(err)
	}
}

func GetGoodss() []Goods {
	var goodss []Goods
	db.Find(&goodss)
	return goodss
}

func GetGoodsById(id int) Goods {
	var goods Goods
	db.First(&goods, id)
	return goods
}
