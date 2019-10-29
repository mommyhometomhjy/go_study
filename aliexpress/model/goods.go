package model

import (
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
}

func FindGoodsByGoodsNo(goodsNo string) Goods {
	var goods Goods
	db.FirstOrCreate(&goods, Goods{GoodsNo: goodsNo})
	return goods
}
func UpdateGoods(goods *Goods) {
	db.Save(goods)
}
