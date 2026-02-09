package utilities

import (
	"math"
	"strconv"
)

func FormatStringToFloat64AndDecimalTo2(stringValue string) float64 {
	value, err := strconv.ParseFloat(stringValue, 64)
	if err != nil {
		return 0.0
	}

	return math.Round(value*100) / 100
}
