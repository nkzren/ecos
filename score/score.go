package score

import (
	"fmt"
	"os"
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

func GetResult(scoreType string) string {
	switch scoreType {
	case "random":
		return valueToLabel(randomScore())
	case "weather":
		return valueToLabel(weatherScore())
	default:
		fmt.Fprintln(os.Stderr, "Invalid score type")
		os.Exit(1)
		return ""
	}
}
