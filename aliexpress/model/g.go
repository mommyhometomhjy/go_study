package model

import (
	"log"

	"aliexpress/config"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

var db *gorm.DB

func SetDB(database *gorm.DB) {
	db = database
}

func ConnectToDB() *gorm.DB {

	driver, dbname := config.GetDataBaseConfig()
	log.Println("Connet to db...")
	db, err := gorm.Open(driver, dbname)
	if err != nil {
		panic("Failed to connect database")
	}
	db.AutoMigrate(&Order{})
	db.AutoMigrate(&Goods{})
	db.AutoMigrate(&OrderDetails{})
	return db
}
