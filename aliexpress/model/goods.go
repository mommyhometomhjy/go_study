package model

import (
	"fmt"

	"math"

	"strings"
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

	//供应商
	GoodsSupplierLink string
}

func CreateGoods(goods *Goods) {
	db.Create(goods)
}

func DeleteGoodsById(id int) {
	var goods Goods
	db.First(&goods, id)

	db.Delete(&goods)
}

func UpdateGoods(goods *Goods) {
	db.Save(goods)
}
func UpdateGoodsPrice() {
	var goodss []Goods
	db.Where("goods_weight >0 and goods_price >0").Find(&goodss)
	for _, goods := range goodss {
		UpdateGoods(&goods)
	}
}

func GetGoodsByGoodsNo(goodsNo string) Goods {
	var goods Goods
	db.FirstOrCreate(&goods, Goods{GoodsNo: strings.ToUpper(goodsNo)})
	return goods
}

func GetGoodss(page int) ([]Goods, BasePage) {
	var total, lastPage, nextPage, currentPage int
	db.Model(&Goods{}).Count(&total)

	//page从1开始
	offset := (page - 1) * 10
	totalPage := int(math.Ceil(float64(total) / 10.0))
	lastPage = page - 1
	currentPage = page
	nextPage = page + 1
	if lastPage < 1 {
		lastPage = 1
	}
	if nextPage > totalPage {
		nextPage = totalPage
	}

	var goodss []Goods
	db.Offset(offset).Limit(10).Find(&goodss)

	basepage := BasePage{
		PrevPage:    lastPage,
		NextPage:    nextPage,
		Total:       totalPage,
		CurrentPage: currentPage,
	}
	return goodss, basepage
}

func GetGoodsById(id int) Goods {
	var goods Goods
	db.First(&goods, id)
	return goods
}

//售价自动计算
func (g *Goods) BeforeSave() {
	if g.GoodsWeight > 0 && g.GoodsPrice > 0 {
		percent, exchange := 0.857, 7.0
		w := fmt.Sprintf("%d", int(math.Ceil(g.GoodsWeight/10.0)*10))

		standShippingCost := GetPriceByWeight(w)
		g.GoodsLastSellPrice = g.GoodsSellPrice
		g.GoodsSellPrice = math.Ceil((g.GoodsPrice+3+standShippingCost)/percent/exchange) - 0.01

	}
}

//创建前把货号字母全部大写
func (g *Goods) BeforeCreate() {
	g.GoodsNo = strings.ToUpper(g.GoodsNo)
}
