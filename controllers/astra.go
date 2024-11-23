package controllers

import (
	"crypto/md5"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/floraorg/sakura/helpers"
	"github.com/gin-gonic/gin"
)

func Astra(c *gin.Context) {
	name := c.Param("name")
	params := helpers.ParseQueryString(c.Request.URL.RawQuery)

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

	hash := md5.Sum([]byte(name))
	points := make([][2]int, 6)
	padding := size / 8

	for i := 0; i < 6; i++ {
		x := int(hash[i*2])*(size-2*padding)/255 + padding
		y := int(hash[i*2+1])*(size-2*padding)/255 + padding
		points[i] = [2]int{x, y}
	}

	var pathData string
	var circles string

	pathData = fmt.Sprintf("M%d,%d ", points[0][0], points[0][1])
	for i := 1; i < len(points); i++ {
		pathData += fmt.Sprintf("L%d,%d ", points[i][0], points[i][1])
	}

	for _, point := range points {
		circles += fmt.Sprintf(`<circle cx="%d" cy="%d" r="%d" fill="white"/>`,
			point[0], point[1], size/40)
	}

	textElement := ""
	gradientOverlay := ""
	if text != "" {
		fontSize := size / 3
		if len(text) > 2 {
			fontSize = size / (len(text) / 2 * 3)
		}

		gradientOverlay = fmt.Sprintf(`
			<defs>
				<radialGradient id="textBackground" cx="50%%" cy="50%%" r="50%%">
					<stop offset="0%%" style="stop-color:#0f0f0f;stop-opacity:0.6"/>
					<stop offset="100%%" style="stop-color:#0f0f0f;stop-opacity:0.6"/>
				</radialGradient>
			</defs>
			<circle cx="%d" cy="%d" r="%d" fill="url(#textBackground)"/>`,
			size, size, size)

		textElement = fmt.Sprintf(`<text x="%d" y="%d" font-size="%d" fill="white" 
			font-family="Arial, Helvetica, sans-serif" font-weight="bold"
			text-anchor="middle" dominant-baseline="middle">%s</text>`,
			size/2, size/2, fontSize, text)
	}

	svg := fmt.Sprintf(`<svg xmlns="http://www.w3.org/2000/svg" width="%d" height="%d" viewBox="0 0 %d %d">
		<rect width="%d" height="%d" fill="#0f0f0f" />
		<path d="%s" stroke="white" stroke-width="%d" fill="none" opacity="0.5"/>
		%s
		%s
		%s
	</svg>`,
		size, size, size, size,
		size, size,
		pathData, size/60,
		circles,
		gradientOverlay,
		textElement)

	c.Header("Content-Type", "image/svg+xml")
	if os.Getenv("ENVIRONMENT") == "DEV" {
		c.Header("Cache-Control", "no-cache")
	} else {
		c.Header("Cache-Control", "public, max-age=604800, immutable")
	}
	c.String(http.StatusOK, svg)
}
