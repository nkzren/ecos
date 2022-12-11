package weather

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/nkzren/ecos/config"
)

type Location struct {
	City    string
	Country string
}

type WeatherbitResponse struct {
	Count int              `json:"count"`
	Data  []WeatherbitData `json:"data"`
}

type WeatherbitData struct {
	Pressure         float64 `json:"pres"`
	Wind_spd         float64 `json:"wind_spd"`
	Temperature      float64 `json:"temp"`
	RelativeHumidity int     `json:"rh"`
	Dhi              float64 `json:"dhi"`
}

var client = &http.Client{Timeout: 5 * time.Second}

func GetData(loc Location) (*WeatherbitData, error) {
	apiKey := config.Root.Weatherbit.Key
	basePath := config.Root.Weatherbit.Path
	resp, err := client.Get(buildUrl(basePath, apiKey, loc.City, loc.Country))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	apiResponse := &WeatherbitResponse{}
	err = json.NewDecoder(resp.Body).Decode(apiResponse)
	if err != nil {
		return nil, err
	}

	if len(apiResponse.Data) == 0 {
		return nil, errors.New("Location not found")
	}

	return &apiResponse.Data[0], nil
}

func buildUrl(basePath, key, city, country string) string {
	return fmt.Sprintf("%s?key=%s&city=%s&country=%s", basePath, key, city, country)
}
