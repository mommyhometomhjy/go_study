package main

import (
	"aliexpress/router"
	"aliexpress/model"
)

func main() {
	db := model.ConnectToDB()
	defer db.Close()
	model.SetDB(db)
	
	router.StartUp()
	
}
