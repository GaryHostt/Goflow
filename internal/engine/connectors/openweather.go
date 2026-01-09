package connectors

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// OpenWeatherAPI fetches weather data from OpenWeather API
type OpenWeatherAPI struct {
	APIKey string
}

// WeatherData represents the weather response structure
type WeatherData struct {
	Weather []struct {
		Main        string `json:"main"`
		Description string `json:"description"`
	} `json:"weather"`
	Main struct {
		Temp      float64 `json:"temp"`
		FeelsLike float64 `json:"feels_like"`
		TempMin   float64 `json:"temp_min"`
		TempMax   float64 `json:"temp_max"`
		Humidity  int     `json:"humidity"`
	} `json:"main"`
	Name string `json:"name"`
}

// FetchWeather retrieves weather data for a city
func (o *OpenWeatherAPI) FetchWeather(city string) Result {
	if o.APIKey == "" {
		return Result{
			Status:  "failed",
			Message: "OpenWeather API key is not configured",
		}
	}

	if city == "" {
		return Result{
			Status:  "failed",
			Message: "City parameter is required",
		}
	}

	url := fmt.Sprintf("https://api.openweathermap.org/data/2.5/weather?q=%s&appid=%s&units=metric", city, o.APIKey)

	resp, err := http.Get(url)
	if err != nil {
		return Result{
			Status:  "failed",
			Message: fmt.Sprintf("Failed to fetch weather data: %v", err),
		}
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return Result{
			Status:  "failed",
			Message: fmt.Sprintf("OpenWeather API returned error status: %d", resp.StatusCode),
		}
	}

	var weatherData WeatherData
	if err := json.NewDecoder(resp.Body).Decode(&weatherData); err != nil {
		return Result{
			Status:  "failed",
			Message: fmt.Sprintf("Failed to parse weather data: %v", err),
		}
	}

	weatherCondition := "Unknown"
	weatherDesc := "Unknown"
	if len(weatherData.Weather) > 0 {
		weatherCondition = weatherData.Weather[0].Main
		weatherDesc = weatherData.Weather[0].Description
	}

	message := fmt.Sprintf("Weather in %s: %s (%s), Temperature: %.1f°C, Humidity: %d%%",
		weatherData.Name,
		weatherCondition,
		weatherDesc,
		weatherData.Main.Temp,
		weatherData.Main.Humidity,
	)

	return Result{
		Status:  "success",
		Message: message,
		Data: map[string]interface{}{
			"city":        weatherData.Name,
			"condition":   weatherCondition,
			"description": weatherDesc,
			"temperature": weatherData.Main.Temp,
			"humidity":    weatherData.Main.Humidity,
		},
	}
}

