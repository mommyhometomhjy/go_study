package router

import (
	"html/template"
	"io"

	"github.com/labstack/echo"
)

var e *echo.Echo

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func init() {

	e = echo.New()

	//静态文件,直接/static访问
	e.Static("/static", "static")

	//设置输出日志
	// e.Use(middleware.Logger())

	//绑定模板
	t := &Template{
		templates: template.Must(template.ParseGlob("views/*/*.html")),
	}
	e.Renderer = t
}
func StartUp() {

	registerRouter()
	e.Logger.Fatal(e.Start(":1323"))
}

func registerRouter() {
	e.GET("/", indexHandler)

	e.GET("/order", orderIndexHandler)
	e.GET("/order/new", orderNewHandler)
	e.POST("/import/order", orderImportExcel)
}
