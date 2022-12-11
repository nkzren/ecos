package score

import (
	"fmt"
	"os"

	"github.com/nkzren/ecoscheduler/weather"
)

// Higher scores are better for scheduling
// 0  <= value < 60  -> bad
// 60 <= value < 100 -> neutral
// value >= 100		 -> good
func valueToLabel(value float64) string {
	if value < 80 {
		return "bad"
	} else if value < 120 {
		return "neutral"
	} else {
		return "good"
	}
}

func GetResult(scoreType string, loc weather.Location) string {
	switch scoreType {
	case "random":
		return valueToLabel(randomScore())
	case "weather":
		result, err := weatherScore(loc)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Could not calculate weather score: %v", err)
		}
		return valueToLabel(result)
	default:
		fmt.Fprintln(os.Stderr, "Invalid score type")
		os.Exit(1)
		return ""
	}
}
