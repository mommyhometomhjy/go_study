package main

import (
	"aliexpress/model"
	"aliexpress/cmd"
)

func main() {
	db := model.ConnectToDB()
	defer db.Close()
	model.SetDB(db)
	
	// cmd.ParseStandShippingCost()
	// cmd.ParseGoodsExcel()
	// cmd.ParseOrderExcel()

	// cmd.ParseShippingCost()
	model.UpdateGoodsPrice()

	cmd.ExportGoodsIncludePrice()
}
