package helpers

import "fmt"

func HashString(str string) int {
	hash := 0
	for _, char := range str {
		hash = int(char) + ((hash << 7) - hash)
	}
	return hash
}

func GenerateUniqueColors(username string) (col1 string, col2 string, background string) {
	hash := HashString(username)
	r := (hash & 0xFF0000) >> 16
	g := (hash & 0x00FF00) >> 8
	b := hash & 0x0000FF

	backgroundR := r / 6
	backgroundG := g / 6
	backgroundB := b / 6

	col1R := ((r * 3) / 2) % 255
	col1G := ((g * 3) / 2) % 255
	col1B := ((b * 3) / 2) % 255

	col2R := (r * 2) % 255
	col2G := (g * 2) % 255
	col2B := (b * 2) % 255

	return fmt.Sprintf("rgb(%d, %d, %d)", backgroundR, backgroundG, backgroundB), fmt.Sprintf("rgb(%d, %d, %d)", col1R, col1G, col1B), fmt.Sprintf("rgb(%d, %d, %d)", col2R, col2G, col2B)
}
