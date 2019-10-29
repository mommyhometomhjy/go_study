package model

import "github.com/jinzhu/gorm"

type StandShippingCost struct {
	gorm.Model
	Weight string
	Price  float64
}

func GetShippingCostByWeight(weight string) StandShippingCost {
	var s StandShippingCost
	db.FirstOrCreate(&s, StandShippingCost{Weight: weight})
	return s
}

func UpdateShippingCost(s *StandShippingCost) {
	db.Save(s)
}
