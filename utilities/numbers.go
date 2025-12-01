package utilities

import (
	"math"
	"math/rand/v2"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

var (
	VietnamesePrinter = message.NewPrinter(language.Vietnamese)
)

func FormatNumber(number int) string {
	formatter := VietnamesePrinter.Sprintf("%d", number)
	return formatter
}

func RoundToInteger(number float64) int {
	return int(math.Round(number))
}

func CeilToInt(f float64) int {
	return int(math.Ceil(f))
}

func FloorToInt(f float64) int {
	return int(math.Floor(f))
}

func RoundByThreshold(n, unit, threshold int) int {
	if unit <= 0 {
		return n // tránh chia 0
	}

	remainder := n % unit
	if remainder >= threshold {
		return n - remainder + unit // làm tròn lên
	}
	return n - remainder // làm tròn xuống
}

// BoolByRatio
// greater than ratio --> true, else false
func BoolByRatio(trueRatio float64) bool {
	value := rand.Float64()
	return value > trueRatio
}
