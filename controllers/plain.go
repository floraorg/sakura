package controllers

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/floraorg/sakura/helpers"
	"github.com/gin-gonic/gin"
)

func Plain(c *gin.Context) {
	name := c.Param("name")

	params := helpers.ParseQueryString(c.Request.URL.RawQuery)

	_, _, c1 := helpers.GenerateUniqueColors(name)

	size := 192
	text := ""

	if val, exists := params["size"]; exists {
		if parsedSize, err := strconv.Atoi(val); err == nil && parsedSize > 0 && parsedSize <= 1000 {
			size = parsedSize
		}
	}
	if val, exists := params["text"]; exists {
		text = val
	}

	textElement := ""
	if text != "" {
		fontSize := size / 3
		if len(text) > 2 {
			fontSize = size / (len(text) / 2 * 3)
		}
		textElement = fmt.Sprintf(`<text x="%d" y="%d" dy="0.1em" text-anchor="middle" dominant-baseline="middle" fill="white" font-family="system-ui, sans-serif" font-weight="bold" font-size="%dpx">%s</text>`,
			size/2, size/2, fontSize, text)
	}

	svg := fmt.Sprintf(`<svg width="%d" height="%d" viewBox="0 0 %d %d" xmlns="http://www.w3.org/2000/svg" style="position:fixed;top:0;left:0;margin:0;padding:0;display:block">
<rect width="%d" height="%d" fill="%s"/>%s</svg>`,
		size, size, size, size,
		size, size, c1,
		textElement)

	c.Header("Content-Type", "image/svg+xml")
	if os.Getenv("ENVIRONMENT") == "DEV" {
		c.Header("Cache-Control", "no-cache")
	} else {
		c.Header("Cache-Control", "public, max-age=604800, immutable")
	}
	c.String(http.StatusOK, svg)
}
