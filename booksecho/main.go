package main

import (
	"booksecho/model"
	"booksecho/router"
)

func main() {
	db := model.ConnectToDB()
	defer db.Close()
	model.SetDB(db)
	router.Start()
}
