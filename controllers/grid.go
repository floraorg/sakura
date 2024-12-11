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

func Grid(c *gin.Context) {
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

	hash := md5.Sum([]byte(name))
	numFilled := (int(hash[0]) % 13) + 8

	gridSize := 6
	cellSize := size / gridSize

	filledPositions := make(map[int]bool)
	hashIndex := 1
	for len(filledPositions) < numFilled && hashIndex < len(hash) {
		pos := int(hash[hashIndex]) % 36
		if !filledPositions[pos] {
			filledPositions[pos] = true
		}
		hashIndex++
	}

	squares := ""
	for i := 0; i < 36; i++ {
		if filledPositions[i] {
			row := i / gridSize
			col := i % gridSize
			x := col * cellSize
			y := row * cellSize
			squares += fmt.Sprintf("<rect x=\"%d\" y=\"%d\" width=\"%d\" height=\"%d\" fill=\"%s\"/>",
				x, y, cellSize, cellSize, c1)
		}
	}

	textElement := ""
	gradientOverlay := ""
	if text != "" {
		fontSize := size / 3
		if len(text) > 2 {
			fontSize = size / (len(text) / 2 * 3)
		}

		gradientOverlay = fmt.Sprintf("<defs><radialGradient id=\"textBackground\" cx=\"50%%\" cy=\"50%%\" r=\"50%%\"><stop offset=\"0%%\" style=\"stop-color:%s;stop-opacity:0.3\"/><stop offset=\"100%%\" style=\"stop-color:%s;stop-opacity:0.3\"/></radialGradient></defs><circle cx=\"%d\" cy=\"%d\" r=\"%d\" fill=\"url(#textBackground)\"/>",
			background, background,
			size/2, size/2, size)

		textElement = fmt.Sprintf("<text x=\"%d\" y=\"%d\" dy=\"0.1em\" font-size=\"%d\" fill=\"%s\" font-weight=\"bold\" font-family=\"Arial, Helvetica, sans-serif\" text-anchor=\"middle\" dominant-baseline=\"middle\">%s</text>",
			size/2, size/2, fontSize, "white", text)
	}

	cornerRadius := 0
	if rounded {
		cornerRadius = size / 8
	}

	svg := fmt.Sprintf("<svg xmlns=\"http://www.w3.org/2000/svg\" width=\"%d\" height=\"%d\" viewBox=\"0 0 %d %d\"><rect width=\"%d\" height=\"%d\" fill=\"%s\" rx=\"%d\" ry=\"%d\"/>%s%s%s</svg>",
		size, size, size, size,
		size, size,
		background,
		cornerRadius, cornerRadius,
		squares,
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
