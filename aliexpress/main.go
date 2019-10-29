package main

import (
	"aliexpress/model"
	"aliexpress/cmd"
)

func main() {
	db := model.ConnectToDB()
	defer db.Close()
	model.SetDB(db)
	cmd.ParseExcel()
}
