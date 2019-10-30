package model

import (
	"fmt"
	"math"
	"strings"

	"github.com/jinzhu/gorm"
)

type Goods struct {
	gorm.Model
	// 产品sku
	GoodsNo string `gorm:"UNIQUE"`

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
		goods.GoodsSellPrice = math.Ceil((goods.GoodsPrice+3+standShippingCost)/percent/exchange) - 0.01
		db.Save(goods)
	}
}

func GetGoodsIncludeSellPriceAndAliexpressId() []Goods {
	var goodss []Goods
	db.Where("goods_sell_price >0 and aliexpress_id <>''").Order("goods_no").Find(&goodss)
	return goodss
}
