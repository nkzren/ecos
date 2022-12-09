package weather

import (
	"math"
)

var r = 0.167226

func GetWindEnergy(data *WeatherbitData) float64 {
	airDensity := getAirDensity(data)

	windSpdCubed := math.Pow(data.Wind_spd, 3)

	return 0.5 * airDensity * windSpdCubed
}

// Values in J/(kg.K)
const dryAirGasConstant = 287.058
const waterVasporGasConstant = 461.495

func getAirDensity(data *WeatherbitData) float64 {
	vaporPressure := func() float64 {
		expoent := func() float64 {
			t := data.Temperature
			return 7.5 * t / (t + 237.3)
		}()
		return 6.1078 * math.Pow(10, expoent) * float64(data.RelativeHumidity)
	}()

	dryPressure := mbarToPa(data.Pressure) - vaporPressure

	return func() float64 {
		kelvin := celsiusToKelvin(data.Temperature)

		getDensity := func(pressure, temperatureInK, gasConstant float64) float64 {
			return pressure / (temperatureInK * gasConstant)
		}

		dryDensity := getDensity(dryPressure, kelvin, dryAirGasConstant)

		vaporDensity := getDensity(vaporPressure, kelvin, waterVasporGasConstant)

		return dryDensity + vaporDensity
	}()
}

func GetSunEnergy(data *WeatherbitData) float64 {
	return data.Dhi
}

func celsiusToKelvin(celsius float64) float64 {
	return celsius + 273
}

func mbarToPa(mbar float64) float64 {
	return mbar * 100
}
