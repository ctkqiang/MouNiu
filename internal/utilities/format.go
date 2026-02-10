package utilities

import (
	"math"
	"strconv"
	"strings"
)

func FormatStringToFloat64AndDecimalTo2(stringValue string) float64 {
	value, err := strconv.ParseFloat(stringValue, 64)
	if err != nil {
		return 0.0
	}

	return math.Round(value*100) / 100
}

func FormatTicker(ticker string) (string, string, error) {
	var (
		Ticker   string
		Exchange string
	)

	parts := strings.Split(ticker, ".")

	if len(parts) == 2 {
		Ticker = parts[0]
		Exchange = parts[1]
	}

	return Ticker, Exchange, nil
}
