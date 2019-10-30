package router

import (
	"aliexpress/vm"
	"html/template"
	"io"
	"net/http"
	"os"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
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
	logF, _ := os.OpenFile("log/log.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "time=${time_rfc3339}, method=${method}, uri=${uri}, status=${status}\r\n",
		Output: logF,
	}))

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
}

func indexHandler(c echo.Context) error {
	vop := vm.IndexViewModelOp{}
	vm := vop.GetVM()

	return c.Render(http.StatusOK, "index", &vm)
}
func orderIndexHandler(c echo.Context) error {
	vop := vm.OrderViewModelOp{}
	vm := vop.GetVM()

	return c.Render(http.StatusOK, "order/index", &vm)
}
