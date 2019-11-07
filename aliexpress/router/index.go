package router

import (
	"github.com/gin-gonic/gin"
)

func indexHandler(c *gin.Context) {

	c.HTML(200, "index", gin.H{
		"Title": "Homepage",
		"Words": "你好",
	})
}
