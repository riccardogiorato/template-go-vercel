package handler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

func MyWeather(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	resp := make(map[string]string)

	// IP User
	userIpAddress := r.RemoteAddr
	// Response API Location
	urlLocationApi := "https://api.freegeoip.app/json/" + userIpAddress + "?apikey=" + os.Getenv("API_KEY_FREEGEOIP")

	fmt.Println("URL Location API:", urlLocationApi)

	userLocationResponse, err := http.Get(urlLocationApi)
	// Get the Location body
	userLocationBody, err := ioutil.ReadAll(userLocationResponse.Body)

	resp["location"] = string(userLocationBody)

	jsonResp, err := json.Marshal(resp)
	if err != nil {
		fmt.Println("Error happened in JSON marshal. Err: %s", err)
	} else {
		w.Write(jsonResp)
	}
	return
}
