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
	r.GET("/plain/:name", controllers.Plain)
	r.GET("/astra/:name", controllers.Astra)
	r.GET("/grid/:name", controllers.Grid)
	r.GET("/colors/:name", controllers.Colors)
	r.GET("/static/:name", controllers.Static)
	r.GET("/floral/:name", controllers.Floral)
	r.GET("/spectra/:name", controllers.Spectra)

	return r
}
