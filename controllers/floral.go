package controllers

import (
	"fmt"
	"math"
	"net/http"
	"os"
	"strconv"

	"github.com/floraorg/sakura/helpers"
	"github.com/gin-gonic/gin"
)

func Floral(c *gin.Context) {
	name := c.Param("name")
	params := helpers.ParseQueryString(c.Request.URL.RawQuery)
	background, c1, c2 := helpers.GenerateUniqueColors(name)

	size := 192
	numPetals := (len(name) % 5) + 12 // This will give us 12-18 petals

	// Parse size parameter
	if val, exists := params["size"]; exists {
		if parsedSize, err := strconv.Atoi(val); err == nil && parsedSize > 0 && parsedSize <= 1000 {
			size = parsedSize
		}
	}

	centerRadius := float64(size) * 0.12 // Smaller center circle (was 0.15)

	petalLengthPercent := 0.30 + float64((int32(helpers.SumASCII(name))%100))/650.0 // Will add 0-0.15
	petalLength := float64(size) * petalLengthPercent

	centerX := float64(size) / 2
	centerY := float64(size) / 2

	// thanks claude i suck at math
	petals := ""
	for i := 0; i < numPetals; i++ {
		angle := (float64(i) * 2 * math.Pi) / float64(numPetals)
		controlPointDist := petalLength * 0.8
		petalWidth := math.Pi / float64(numPetals) * 1.2

		x1 := centerX + math.Cos(angle)*centerRadius
		y1 := centerY + math.Sin(angle)*centerRadius
		x2 := centerX + math.Cos(angle)*petalLength
		y2 := centerY + math.Sin(angle)*petalLength

		// Control points for the bezier curve
		cx1 := centerX + math.Cos(angle-petalWidth)*controlPointDist
		cy1 := centerY + math.Sin(angle-petalWidth)*controlPointDist
		cx2 := centerX + math.Cos(angle+petalWidth)*controlPointDist
		cy2 := centerY + math.Sin(angle+petalWidth)*controlPointDist

		// Create a petal using a closed bezier curve
		petals += fmt.Sprintf(`<path d="M %.2f %.2f Q %.2f %.2f %.2f %.2f Q %.2f %.2f %.2f %.2f Z" fill="%s"/>`,
			x1, y1,
			cx1, cy1,
			x2, y2,
			cx2, cy2,
			x1, y1,
			c2)
	}

	text := ""
	textElement := ""
	gradientOverlay := ""
	if val, exists := params["text"]; exists {
		text = val
	}
	if text != "" {
		fontSize := size / 3
		if len(text) > 2 {
			fontSize = size / (len(text) / 2 * 4)
		}
		gradientOverlay = fmt.Sprintf("<defs><radialGradient id=\"textBackground\" cx=\"50%%\" cy=\"50%%\" r=\"50%%\"><stop offset=\"0%%\" style=\"stop-color:%s;stop-opacity:0.5\"/><stop offset=\"100%%\" style=\"stop-color:%s;stop-opacity:0.5\"/></radialGradient></defs><circle cx=\"%d\" cy=\"%d\" r=\"%d\" fill=\"url(#textBackground)\"/>",
			background, background,
			size/2, size/2, size)
		textElement = fmt.Sprintf(`<text x="%d" y="%d" text-anchor="middle" dominant-baseline="middle" fill="white" font-family="system-ui, sans-serif" font-weight="bold" font-size="%dpx">%s</text>`,
			size/2, size/2, fontSize, text)
	}

	svg := fmt.Sprintf(`<svg width="%d" height="%d" viewBox="0 0 %d %d" xmlns="http://www.w3.org/2000/svg">
<rect width="%d" height="%d" fill="%s"/>
%s
<circle cx="%f" cy="%f" r="%f" stroke="white" stroke-width="%f" fill="%s"/>
%s %s
</svg>`,
		size, size, size, size,
		size, size, background,
		petals,
		centerX, centerY, centerRadius, centerRadius*0.1, c1, gradientOverlay,
		textElement)

	c.Header("Content-Type", "image/svg+xml")
	if os.Getenv("ENVIRONMENT") == "DEV" {
		c.Header("Cache-Control", "no-cache")
	} else {
		c.Header("Cache-Control", "public, max-age=604800, immutable")
	}
	c.String(http.StatusOK, svg)
}
