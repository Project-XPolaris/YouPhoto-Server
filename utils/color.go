package utils

import (
	"math"
	"strconv"
)

func HexToRGB(hex string) (r, g, b float64) {
	hexValue, _ := strconv.ParseInt(hex[1:], 16, 32)
	r = float64(hexValue >> 16 & 0xFF)
	g = float64(hexValue >> 8 & 0xFF)
	b = float64(hexValue & 0xFF)
	return
}

func ColorDistance(color1, color2 string) float64 {
	r1, g1, b1 := HexToRGB(color1)
	r2, g2, b2 := HexToRGB(color2)
	distance := math.Sqrt(math.Pow(r2-r1, 2) + math.Pow(g2-g1, 2) + math.Pow(b2-b1, 2))
	return distance
}
