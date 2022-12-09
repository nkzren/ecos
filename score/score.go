package score

import (
	"fmt"
	"os"

	"github.com/nkzren/ecoscheduler/weather"
)

// Takes a value between 0 and 100 and returns a label
// Higher scores are better for scheduling
// 0  <= value < 40   -> bad
// 40 <= value < 60   -> neutral
// 60 <= value < 100  -> good
func valueToLabel(value int) string {
	if value < 40 {
		return "bad"
	} else if value < 60 {
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
