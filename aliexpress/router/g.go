package router

import (
	"github.com/gin-gonic/gin"
)

var e *gin.Engine

func init() {

	e = gin.Default()

	//静态文件,直接/static访问
	e.Static("/static", "static")
	e.Use(gin.Logger())
	//设置输出日志
	// e.Use(middleware.Logger())
	e.LoadHTMLGlob("views/*/*.html")
	//绑定模板

}
func StartUp() {

	registerRouter()
	e.Run(":1323")
}

func registerRouter() {

	index := e.Group("/index")
	{
		index.GET("/", indexHandler)
	}

	order := e.Group("/order")
	{
		//获取全部order
		order.GET("/index", orderIndexHandler)

		//创建Order页面
		order.GET("/new", orderNewHandler)
		//创建order
		order.POST("/new", orderCreate)

		//删除order
		order.POST("/delete/:id", orderDelete)

		//获取指定order
		order.GET("/edit/:id", orderEditHandler)
		//修改指定order
		order.POST("/edit/:id", orderUpdate)

		//批量导入order
		order.POST("/import", orderImportExcel)
	}

	goods := e.Group("/goods")
	{

		goods.GET("/index", getGoodss)

		//创建goods页面
		goods.GET("/new", goodsNew)
		//创建goods
		goods.POST("/new", goodsCreate)

		//删除goods
		goods.POST("/delete/:id", goodsDelete)

		//获取指定goods
		goods.GET("/edit/:id", goodsEdit)
		// //修改指定goods
		goods.POST("/edit/:id", goodsUpdate)

		//导入标准运费
		goods.POST("/importstandardshippingcost", parseStandardShippingCost)

		//导出价格有变动的
		goods.GET("/exportsellpricechanged", exportsellpricechanged)
	}
}

type BasePageViewModel struct {
	PrevPage    int
	NextPage    int
	Total       int
	CurrentPage int
	Limit       int
}

// SetPrevAndNextPage func
func (v *BasePageViewModel) SetPrevAndNextPage() {
	if v.CurrentPage > 1 {
		v.PrevPage = v.CurrentPage - 1
	}

	if (v.Total-1)/v.Limit >= v.CurrentPage {
		v.NextPage = v.CurrentPage + 1
	}
}

// SetBasePageViewModel func
func (v *BasePageViewModel) SetBasePageViewModel(total, page, limit int) {
	v.Total = total
	v.CurrentPage = page
	v.Limit = limit
	v.SetPrevAndNextPage()
}
