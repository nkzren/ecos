package score

import "github.com/nkzren/ecoscheduler/weather"

func weatherScore(loc weather.Location) (float64, error) {
	data, err := weather.GetData(loc)
	if err != nil {
		return -1, err
	}
	return calculateScore(data), nil
}

func calculateScore(data *weather.WeatherbitData) float64 {
	return weather.GetWindEnergy(data) + weather.GetSunEnergy(data)
}
