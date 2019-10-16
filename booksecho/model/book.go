package model

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Book struct {
	ID       int    `gorm:"primary_key";json:"id";from:"id"`
	Title    string `gorm:"type:varchar(64)";json:"title";from:"title"`
	Subtitle string `gorm:"type:varchar(64)";json:"subtitle";from:"subtitle"`
	Pic      string `gorm:"type:varchar(200)";json:"pic";from:"pic"`
	Author   string `gorm:"type:varchar(64)";json:"author";from:"author"`
	Summary  string `gorm:"type:varchar(500)";json:"summary";from:"summary"`
	Isbn     string `gorm:"type:varchar(20);index:isbn";json:"isbn";from:"isbn"`
}

func GetBookByIsbn(isbn string) (*Book, error) {
	book := Book{Isbn: isbn}
	err := db.Where("isbn = ?", book.Isbn).First(&book).Error

	if err != nil {
		db.Create(&book)
		fetchBook(&book)
	}
	return &book, nil

}
func fetchBook(book *Book) {
	cliend := &http.Client{}

	req, _ := http.NewRequest("GET", "https://jisuisbn.market.alicloudapi.com/isbn/query?isbn="+book.Isbn, nil)

	req.Header.Add("Authorization", "APPCODE b8ba9d0f2f294329850e3956048c1e32")

	resp, err := cliend.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()

	result := new(Result)

	pix, _ := ioutil.ReadAll(resp.Body)

	json.Unmarshal(pix, result)

	// fmt.Println(result.Status, result.Msg)
	if result.Status == "0" {

		book.Author = result.Book.Author

		book.Pic = result.Book.Pic

		book.Subtitle = result.Book.Subtitle

		book.Summary = result.Book.Summary

		book.Title = result.Book.Title
	} else {
		book.Title = "查不到此书,请手动编辑"
	}
	db.Save(book)

}

func GetAllBooks() []Book {
	var books []Book
	db.Find(&books)
	return books
}

func GetBookById(id int) (*Book, error) {

	book := Book{ID: id}
	err := db.First(&book, id).Error
	if err != nil {
		return nil, err
	}
	return &book, nil
}

func DeleteBookById(id int) {
	book := Book{ID: id}
	db.Delete(&book)
}
func UpdateBook(book *Book) {
	db.Save(book)
}
func CreateBook(book *Book) {
	db.Create(&book)
}
