package controllers

import (
	"github.com/floraorg/sakura/helpers"
	"github.com/gin-gonic/gin"
)

func Colors(c *gin.Context) {
	name := c.Param("name")

	bg, c2, c1 := helpers.GenerateUniqueColors(name)
	c.JSON(200, gin.H{
		"primary":    c1,
		"secondary":  c2,
		"background": bg,
	})

}
