package routes

import (
	"github.com/floraorg/sakura/controllers"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()
	r.LoadHTMLGlob("views/*")
	r.GET("/", controllers.Index)

	r.GET("/linear/:name", controllers.Linear)

	return r
}
