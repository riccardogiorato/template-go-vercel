package handler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

type User struct {
	ID    int32  `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func MyWeather(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	resp := make(map[string]string)

	// IP User
	userIpAddress := "93.34.228.207" // r.RemoteAddr
	// Response API Location
	urlLocationApi := "https://api.freegeoip.app/json/" + userIpAddress + "?apikey=" + os.Getenv("API_KEY_FREEGEOIP")
	fmt.Println("Location API URL:", urlLocationApi)
	userLocationResponse, err := http.Get(urlLocationApi)
	// Get the Location body
	userLocationBody, err := ioutil.ReadAll(userLocationResponse.Body)
	userLocationJson := string(userLocationBody)
	fmt.Println("Location API Json String:", userLocationJson)
	var userLocation map[string]interface{}
	json.Unmarshal([]byte(userLocationBody), &userLocation)

	resp["latitude"] = fmt.Sprint(userLocation["latitude"].(float64))
	resp["longitude"] = fmt.Sprint(userLocation["longitude"].(float64))

	// Response API Weather
	urlWeatherApi := "https://api.openweathermap.org/data/2.5/weather?lat=" + resp["latitude"] + "&lon=" + resp["longitude"] + "&appid=" + os.Getenv("API_KEY_OPENWEATHER")
	fmt.Println("Weather API Json String:", urlWeatherApi)
	weatherApiResponse, err := http.Get(urlWeatherApi)
	weatherApiBody, err := ioutil.ReadAll(weatherApiResponse.Body)
	weatherApiJson := string(weatherApiBody)

	resp["weather"] = weatherApiJson

	jsonResp, err := json.Marshal(resp)
	if err != nil {
		fmt.Println("Error happened in JSON marshal. Err:", err)
	} else {
		w.Write(jsonResp)
	}
	return
}
