package util

import (
	"math"
	"strings"
)

func MathRoundPrecision(value float64, precision int) float64 {
	return math.Round(value*(math.Pow10(precision))) / math.Pow10(precision)
}

func FormatTitle(value string) string {
	return strings.ToUpper(strings.Join(strings.Fields(value), " "))
}

func FormatTextWithoutSpace(value string) string {
	return strings.Join(strings.Fields(value), "")
}

func IsAlpha(value string) bool {
	for _, r := range value {
		if (r < 'a' || r > 'z') && (r < 'A' || r > 'Z') {
			return false
		}
	}
	return true
}
