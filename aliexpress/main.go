package main

import (
	"aliexpress/model"
)

func main() {
	db := model.ConnectToDB()
	defer db.Close()
	model.SetDB(db)

}
