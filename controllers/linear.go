package controllers

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/floraorg/sakura/helpers"
	"github.com/gin-gonic/gin"
)

func Linear(c *gin.Context) {
	name := c.Param("name")

	// Get the raw query string
	rawQuery := c.Request.URL.RawQuery

	// Manual query parsing
	size := 192
	text := ""

	if rawQuery != "" {
		params := strings.Split(rawQuery, "?")
		for _, param := range params {
			if strings.Contains(param, "=") {
				parts := strings.Split(param, "=")
				key := parts[0]
				value := parts[1]

				switch key {
				case "size":
					if parsedSize, err := strconv.Atoi(value); err == nil && parsedSize > 0 && parsedSize <= 1000 {
						size = parsedSize
					}
				case "text":
					text = value
				}
			}
		}
	}

	// Generate colors
	_, c1, c2 := helpers.GenerateUniqueColors(name)

	// Create text element if needed
	textElement := ""
	if text != "" {
		fontSize := size / 3
		if len(text) > 2 {
			fontSize = size / (len(text)/2 + 1)
		}
		textElement = fmt.Sprintf(`<text x="%d" y="%d" text-anchor="middle" dominant-baseline="middle" fill="white" font-family="system-ui, sans-serif" font-weight="bold" font-size="%dpx">%s</text>`,
			size/2, size/2, fontSize, text)
	}

	svg := fmt.Sprintf(`<svg width="%d" height="%d" viewBox="0 0 %d %d" xmlns="http://www.w3.org/2000/svg" style="position:fixed;top:0;left:0;margin:0;padding:0;display:block">
<defs>
    <linearGradient id="g" x1="0" y1="0" x2="1" y2="1">
        <stop offset="0%%" stop-color="%s"/>
        <stop offset="100%%" stop-color="%s"/>
    </linearGradient>
</defs>
<rect width="%d" height="%d" fill="url(#g)"/>%s</svg>`,
		size, size, size, size,
		c1, c2,
		size, size,
		textElement)

	c.Header("Content-Type", "image/svg+xml")
	if os.Getenv("ENVIRONMENT") == "DEV" {
		c.Header("Cache-Control", "no-cache")
	} else {
		c.Header("Cache-Control", "public, max-age=604800, immutable")
	}
	c.String(http.StatusOK, svg)
}
