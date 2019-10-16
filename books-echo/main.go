package main

import (
	"books/model"

	"strconv"

	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	db := model.ConnectToDB()
	defer db.Close()
	model.SetDB(db)

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	e.GET("/books", getAllBooks)

	e.GET("/book/isbn/:isbn", getBookByIsbn)

	e.GET("/book/:id", getBookById)
	e.DELETE("/book/:id", deleteBookById)
	e.PUT("/book/:id", updateBook)

	e.Start(":1323")
}

func getBookByIsbn(c echo.Context) error {
	isbn := c.Param("isbn")
	book, _ := model.GetBookByIsbn(isbn)
	return c.JSON(200, *book)

}
func getAllBooks(c echo.Context) error {
	return c.JSON(200, model.GetAllBooks())

}
func getBookById(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	book, _ := model.GetBookById(id)
	return c.JSON(200, book)

}
func deleteBookById(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	model.DeleteBookById(id)
	return c.String(200, "删除成功")

}
func updateBook(c echo.Context) error {

	id, _ := strconv.Atoi(c.FormValue("id"))
	book, _ := model.GetBookById(id)
	book.Isbn = c.FormValue("isbn")
	book.Author = c.FormValue("author")
	book.Subtitle = c.FormValue("subtitle")
	book.Summary = c.FormValue("summary")
	book.Title = c.FormValue("title")

	model.UpdateBook(book)
	return c.JSON(200, book)

}
