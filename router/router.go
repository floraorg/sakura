package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/floraorg/sakura/controllers"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()
	r.GET("/", controllers.Index)

	return r
}