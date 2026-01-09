package connectors

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// OpenWeatherAPI handles OpenWeather API integrations
type OpenWeatherAPI struct {
	APIKey string
}

// WeatherData represents the OpenWeather API response
type WeatherData struct {
	Main struct {
		Temp     float64 `json:"temp"`
		Humidity int     `json:"humidity"`
	} `json:"main"`
	Weather []struct {
		Main        string `json:"main"`
		Description string `json:"description"`
	} `json:"weather"`
	Name string `json:"name"`
}

// FetchWeather retrieves weather data for a city
func (w *OpenWeatherAPI) FetchWeather(city string) Result {
	return w.FetchWeatherWithContext(context.Background(), city)
}

// FetchWeatherWithContext retrieves weather data with context awareness
func (w *OpenWeatherAPI) FetchWeatherWithContext(ctx context.Context, city string) Result {
	start := time.Now()

	select {
	case <-ctx.Done():
		return NewCancelledResult("Context cancelled before weather request: " + ctx.Err().Error())
	default:
	}

	url := fmt.Sprintf("https://api.openweathermap.org/data/2.5/weather?q=%s&appid=%s&units=metric", city, w.APIKey)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return NewFailureResult(fmt.Sprintf("Failed to create weather request: %v", err), start)
	}

	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	resp, err := client.Do(req)

	select {
	case <-ctx.Done():
		return NewCancelledResult("Context cancelled during weather request: " + ctx.Err().Error())
	default:
	}

	if err != nil {
		return NewFailureResult(fmt.Sprintf("OpenWeather API request failed: %v", err), start)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return NewFailureResult(fmt.Sprintf("OpenWeather returned error status: %d", resp.StatusCode), start)
	}

	var weather WeatherData
	if err := json.NewDecoder(resp.Body).Decode(&weather); err != nil {
		return NewFailureResult(fmt.Sprintf("Failed to decode weather response: %v", err), start)
	}

	description := "N/A"
	if len(weather.Weather) > 0 {
		description = weather.Weather[0].Description
	}

	return NewSuccessResult(fmt.Sprintf("Weather in %s: %s, %.1f°C", city, description, weather.Main.Temp), map[string]interface{}{
		"city":        city,
		"temperature": weather.Main.Temp,
		"humidity":    weather.Main.Humidity,
		"description": description,
	}, start)
}
