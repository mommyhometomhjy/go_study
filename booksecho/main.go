package main

import (
	"booksecho/model"

	"strconv"

	"html/template"
	"io"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Template struct {
	templates *template.Template
}
type IndexVeiws struct {
	Title string
	Data  interface{}
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func main() {
	db := model.ConnectToDB()
	defer db.Close()
	model.SetDB(db)

	e := echo.New()
	e.Static("/static", "public/static")
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	t := &Template{
		templates: template.Must(template.ParseGlob("public/views/**/*.html")),
	}
	e.Renderer = t

	//获取所有的书籍
	e.GET("/books", getAllBooks)

	//打开新建书籍页面
	e.GET("/book/new", newBook)

	//创建新的书籍
	e.POST("/books", createBook)

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
	v := IndexVeiws{Title: "所有书籍", Data: model.GetAllBooks()}
	return c.Render(http.StatusOK, "books/index", &v)

}

func newBook(c echo.Context) error {
	v := IndexVeiws{Title: "新建书籍"}
	return c.Render(http.StatusOK, "books/new", &v)
}
func createBook(c echo.Context) (err error) {
	book := new(model.Book)

	if err = c.Bind(book); err != nil {
		return c.String(200, "创建失败")
	}

	model.CreateBook(book)
	return c.String(200, "创建成功")
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
