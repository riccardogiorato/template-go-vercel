package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type locationResponse struct {
	Success   bool    `json:"success"`
	Message   string  `json:"message"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	City      string  `json:"city"`
	Country   string  `json:"country"`
}

type weatherResponse struct {
	Current struct {
		Temperature float64 `json:"temperature_2m"`
		WeatherCode int     `json:"weather_code"`
	} `json:"current"`
}

func MyWeather(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userIPAddress := "93.34.228.207" // r.RemoteAddr includes a port in local development.
	locationURL := "https://ipwho.is/" + userIPAddress
	var location locationResponse
	if err := getJSON(locationURL, &location); err != nil {
		http.Error(w, fmt.Sprintf("Error getting location: %v", err), http.StatusBadGateway)
		return
	}
	if !location.Success {
		http.Error(w, fmt.Sprintf("Error getting location: %s", location.Message), http.StatusBadGateway)
		return
	}

	weatherURL := fmt.Sprintf(
		"https://api.open-meteo.com/v1/forecast?latitude=%f&longitude=%f&current=temperature_2m,weather_code",
		location.Latitude,
		location.Longitude,
	)
	var weather weatherResponse
	if err := getJSON(weatherURL, &weather); err != nil {
		http.Error(w, fmt.Sprintf("Error getting weather: %v", err), http.StatusBadGateway)
		return
	}

	w.WriteHeader(http.StatusCreated)
	err := json.NewEncoder(w).Encode(map[string]interface{}{
		"city":         location.City,
		"country":      location.Country,
		"latitude":     location.Latitude,
		"longitude":    location.Longitude,
		"temperature":  weather.Current.Temperature,
		"weather_code": weather.Current.WeatherCode,
		"github":       "https://github.com/riccardogiorato/template-go-vercel/blob/main/api/myweather.go",
	})
	if err != nil {
		fmt.Println("Error happened in JSON marshal. Err:", err)
	}
}

func getJSON(url string, target interface{}) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return fmt.Errorf("HTTP %d", resp.StatusCode)
	}
	return json.NewDecoder(resp.Body).Decode(target)
}
