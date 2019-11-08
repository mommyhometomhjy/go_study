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

func GetGoodsById(id int) Goods {
	var goods Goods
	db.First(&goods, id)
	return goods
}

func SearchGoodsByGoodsNo(no string, page, limit int) ([]Goods, int) {
	var total int
	db.Model(&Goods{}).Where("goods_no like ?", "%"+no+"%").Count(&total)

	offset := (page - 1) * limit

	var goodss []Goods
	db.Offset(offset).Limit(limit).Where("goods_no like ?", "%"+no+"%").Find(&goodss)
	return goodss, total
}

func GetGoodsSellPriceChanged() []Goods {
	var goodss []Goods
	db.Where("goods_sell_price <> goods_last_sell_price").Find(&goodss)
	return goodss
}

//售价自动计算
func (g *Goods) BeforeSave() {
	if g.GoodsWeight > 0 && g.GoodsPrice > 0 {
		w := fmt.Sprintf("%d", int(math.Ceil(g.GoodsWeight/10.0)*10))

		standShippingCost := GetPriceByWeight(w)
		g.GoodsLastSellPrice = g.GoodsSellPrice
		g.GoodsSellPrice = math.Ceil((g.GoodsPrice+PROFIT+standShippingCost)/PERCENT/EXCHANGE) - 0.01

	}
}

//创建前把货号字母全部大写
func (g *Goods) BeforeCreate() {
	g.GoodsNo = strings.ToUpper(g.GoodsNo)
}
