package score

import (
	"fmt"
	"os"

	"github.com/nkzren/ecoscheduler/weather"
)

func weatherScore(loc weather.Location) (float64, error) {
	data, err := weather.GetData(loc)
	if err != nil {
		return -1, err
	}
	score := calculateScore(data)
	logScore(loc, score)
	return score, nil
}

func calculateScore(data *weather.WeatherbitData) float64 {
	return weather.GetWindEnergy(data) + weather.GetSunEnergy(data)
}

func logScore(loc weather.Location, score float64) {
	fmt.Fprintf(os.Stdout, "Score for location (%s, %s): %.2f\n", loc.City, loc.Country, score)
}
