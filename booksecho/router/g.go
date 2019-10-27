package router

import (
	"booksecho/model"
	"os"
	"path"
	"path/filepath"
	"strconv"

	"html/template"
	"io"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Template struct {
	templates map[string]*template.Template
}
type IndexVeiws struct {
	Title    string
	Disabled string
	Data     interface{}
}

var (
	e         *echo.Echo
	templates map[string]*template.Template
)

func init() {
	e = echo.New()
	e.Static("/static", "public/static")
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	t := &Template{
		templates: PopulateTemplates(),
	}
	e.Renderer = t

	registerRouters()

}
func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {

	return t.templates[name].Execute(w, data)
}

func Start() {
	e.Start(":1323")
}

func registerRouters() {
	//获取所有的书籍
	e.GET("/books", getAllBooks)

	//打开新建书籍页面
	e.GET("/book/new", newBook)

	//创建新的书籍
	e.POST("/books", createBook)

	//展示书籍
	e.GET("/book/:id", getBookById)

	//编辑书籍
	e.GET("/book/:id/edit", editBookById)

	e.DELETE("/book/:id", deleteBookById)

	e.POST("/book/:id", updateBook)

	e.GET("/book/isbn/:isbn", getBookByIsbn)
}

func getAllBooks(c echo.Context) error {
	v := IndexVeiws{Title: "所有书籍", Data: model.GetAllBooks()}
	return c.Render(http.StatusOK, "index.html", &v)
}

func newBook(c echo.Context) error {
	book := new(model.Book)
	v := IndexVeiws{Title: "新建书籍", Data: book}
	return c.Render(http.StatusOK, "new.html", &v)
}
func createBook(c echo.Context) (err error) {
	book := new(model.Book)

	if err = c.Bind(book); err != nil {
		return c.String(200, "创建失败")
	}
	file, err := c.FormFile("file")
	if err == nil {
		src, err := file.Open()
		if err != nil {
			return err
		}
		defer src.Close()

		dir, _ := os.Getwd()

		dst, err := os.Create(filepath.Join(dir, "public", "static", "images", book.Isbn+path.Ext(file.Filename)))

		if err != nil {
			return err
		}

		defer dst.Close()
		if _, err = io.Copy(dst, src); err != nil {
			return err
		}
		book.Pic = "/static/images/" + book.Isbn + path.Ext(file.Filename)
	}

	model.CreateBook(book)
	v := IndexVeiws{Title: "查看书籍", Data: book}
	return c.Render(200, "show.html", &v)
}
func getBookById(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	book, _ := model.GetBookById(id)
	v := IndexVeiws{Title: "查看书籍", Data: book}
	return c.Render(200, "show.html", &v)

}

func editBookById(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	book, _ := model.GetBookById(id)
	v := IndexVeiws{Title: "编辑书籍", Data: book}
	return c.Render(200, "edit.html", &v)
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

	file, err := c.FormFile("file")
	if err == nil {
		src, err := file.Open()
		if err != nil {
			return err
		}
		defer src.Close()

		dir, _ := os.Getwd()

		dst, err := os.Create(filepath.Join(dir, "public", "static", "images", book.Isbn+path.Ext(file.Filename)))

		if err != nil {
			return err
		}

		defer dst.Close()
		if _, err = io.Copy(dst, src); err != nil {
			return err
		}
		book.Pic = "/static/images/" + book.Isbn + path.Ext(file.Filename)
	}

	model.UpdateBook(book)
	v := IndexVeiws{Title: "查看书籍", Data: book}
	return c.Render(200, "show.html", &v)
}

func getBookByIsbn(c echo.Context) error {
	isbn := c.Param("isbn")
	book, _ := model.GetBookByIsbn(isbn)
	// v := IndexVeiws{Title: "编辑书籍", Data: book}
	return c.JSON(200, book.ID)
}
