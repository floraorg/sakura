package controllers

import (
	"crypto/md5"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"strconv"

	"github.com/floraorg/sakura/helpers"
	"github.com/gin-gonic/gin"
)

func Static(c *gin.Context) {
	name := c.Param("name")
	params := helpers.ParseQueryString(c.Request.URL.RawQuery)

	// Get colors from helper
	background, _, c1 := helpers.GenerateUniqueColors(name)

	size := 192
	rounded := false
	text := ""

	if val, exists := params["size"]; exists {
		if parsedSize, err := strconv.Atoi(val); err == nil && parsedSize > 0 && parsedSize <= 1000 {
			size = parsedSize
		}
	}
	if val, exists := params["text"]; exists {
		text = val
	}
	if _, exists := params["rounded"]; exists {
		rounded = true
	}

	// Use hash to seed the random number generator
	hash := md5.Sum([]byte(name))
	seed := int64(hash[0]) | int64(hash[1])<<8 | int64(hash[2])<<16 | int64(hash[3])<<24
	rand.Seed(seed)

	// Generate static pixels
	pixelSize := size / 32 // Controls the "graininess" of the static
	numPixels := 32 * 32   // Total number of pixels

	pixels := ""
	for i := 0; i < numPixels; i++ {
		x := (i % 32) * pixelSize
		y := (i / 32) * pixelSize

		// Randomly decide pixel opacity
		opacity := rand.Float64()*0.8 + 0.1 // opacity between 0.1 and 0.9

		pixels += fmt.Sprintf("<rect x=\"%d\" y=\"%d\" width=\"%d\" height=\"%d\" fill=\"%s\" opacity=\"%.2f\"/>",
			x, y, pixelSize, pixelSize, c1, opacity)
	}

	// Create text overlay if specified
	textElement := ""
	gradientOverlay := ""
	if text != "" {
		fontSize := size / 3
		if len(text) > 2 {
			fontSize = size / (len(text) / 2 * 3)
		}

		gradientOverlay = fmt.Sprintf("<defs><radialGradient id=\"textBackground\" cx=\"50%%\" cy=\"50%%\" r=\"50%%\"><stop offset=\"0%%\" style=\"stop-color:%s;stop-opacity:0.5\"/><stop offset=\"100%%\" style=\"stop-color:%s;stop-opacity:0.5\"/></radialGradient></defs><circle cx=\"%d\" cy=\"%d\" r=\"%d\" fill=\"url(#textBackground)\"/>",
			background, background,
			size/2, size/2, size)

		textElement = fmt.Sprintf("<text x=\"%d\" y=\"%d\" font-size=\"%d\" font-weight=\"bold\" fill=\"%s\" font-family=\"Arial, Helvetica, sans-serif\" text-anchor=\"middle\" dominant-baseline=\"middle\">%s</text>",
			size/2, size/2, fontSize, "white", text)
	}

	cornerRadius := 0
	if rounded {
		cornerRadius = size / 8
	}

	// Generate the final SVG
	svg := fmt.Sprintf("<svg xmlns=\"http://www.w3.org/2000/svg\" width=\"%d\" height=\"%d\" viewBox=\"0 0 %d %d\"><rect width=\"%d\" height=\"%d\" fill=\"%s\" rx=\"%d\" ry=\"%d\"/>%s%s%s</svg>",
		size, size, size, size,
		size, size,
		background,
		cornerRadius, cornerRadius,
		pixels,
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
