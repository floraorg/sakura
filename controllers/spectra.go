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

func Spectra(c *gin.Context) {
	name := c.Param("name")
	params := helpers.ParseQueryString(c.Request.URL.RawQuery)
	background, _, c2 := helpers.GenerateUniqueColors(name)

	size := 192
	if val, exists := params["size"]; exists {
		if parsedSize, err := strconv.Atoi(val); err == nil && parsedSize > 0 && parsedSize <= 1000 {
			size = parsedSize
		}
	}

	numWaves := (len(name) % 3) + 3 // 3-5 waves
	nameSum := helpers.SumASCII(name)

	waves := ""
	for i := 0; i < numWaves; i++ {
		amplitude := float64(size) * (0.1 + float64(i)*0.05) // Increasing amplitudes
		frequency := 2 + float64((nameSum+i)%3)              // Varies frequency based on name
		phase := float64(nameSum%10) * 0.1 * float64(i)      // Varying phase shift

		// Generate wave path
		path := "M"
		for x := 0; x <= size; x += 2 {
			// Calculate y position with sine wave
			xPos := float64(x)
			yPos := float64(size)/2 +
				amplitude*math.Sin(2*math.Pi*frequency*xPos/float64(size)+phase)

			if x == 0 {
				path += fmt.Sprintf(" %.2f,%.2f", xPos, yPos)
			} else {
				path += fmt.Sprintf(" L %.2f,%.2f", xPos, yPos)
			}
		}

		waves += fmt.Sprintf(`<path d="%s" fill="none" stroke="%s" stroke-width="%.1f" 
			stroke-linecap="round" opacity="%.2f"/>`,
			path, c2, float64(size)*0.005,
			0.3+float64(i)*0.2) // Increasing opacity for each wave
	}

	// Handle text overlay
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
		gradientOverlay = fmt.Sprintf(`<defs>
			<radialGradient id="textBackground" cx="50%%" cy="50%%" r="50%%">
				<stop offset="0%%" style="stop-color:%s;stop-opacity:0.5"/>
				<stop offset="100%%" style="stop-color:%s;stop-opacity:0.5"/>
			</radialGradient>
		</defs>
		<circle cx="%d" cy="%d" r="%d" fill="url(#textBackground)"/>`,
			background, background,
			size/2, size/2, size/3)

		textElement = fmt.Sprintf(`<text x="%d" y="%d" 
			text-anchor="middle" dominant-baseline="middle" 
			fill="white" font-family="system-ui, sans-serif" 
			font-weight="bold" font-size="%dpx">%s</text>`,
			size/2, size/2, fontSize, text)
	}

	svg := fmt.Sprintf(`<svg width="%d" height="%d" viewBox="0 0 %d %d" 
		xmlns="http://www.w3.org/2000/svg">
		<rect width="%d" height="%d" fill="%s"/>
		%s
		%s
		%s
	</svg>`,
		size, size, size, size,
		size, size, background,
		waves,
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
