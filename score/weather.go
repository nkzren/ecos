package score

import "github.com/nkzren/ecoscheduler/weather"

func weatherScore(loc weather.Location) (int, error) {
	data, err := weather.GetData(loc)
	if err != nil {
		return -1, err
	}
	return calculateScore(data), nil
}

func calculateScore(data *weather.WeatherbitData) int {
	return int(weather.GetWindEnergy(data) + weather.GetSunEnergy(data))
}
