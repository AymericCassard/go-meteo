package main

import (
	"encoding/json"
	"fmt"
	"github.com/a-h/templ"
	"go-meteo/controller"
	"io"
	"net/http"
	"strconv"
)

type ForecastResponse struct {
	Latitude             float64 `json:"latitude"`
	Longitude            float64 `json:"longitude"`
	GenerationTimeMs     float64 `json:"generationtime_ms"`
	UTCOffsetSeconds     int     `json:"utc_offset_seconds"`
	Timezone             string  `json:"timezone"`
	TimezoneAbbreviation string  `json:"timezone_abbreviation"`
	Elevation            float64 `json:"elevation"`
	CurrentUnits         struct {
		Time          string `json:"time"`
		Interval      string `json:"interval"`
		Temperature2m string `json:"temperature_2m"`
	} `json:"current_units"`
	Current struct {
		Time          string  `json:"time"`
		Interval      int     `json:"interval"`
		Temperature2m float64 `json:"temperature_2m"`
	} `json:"current"`
}

func main() {
	forecast, err := API_test()
	if err != nil {
		fmt.Printf("Erreur : %s\n", err)
		return
	}
	temp := strconv.Itoa(int(forecast.Current.Temperature2m))
	index := index(controller.Default(temp))
	http.Handle("/", templ.Handler(index))
	fmt.Println("Listening on http://localhost:3000")
	http.ListenAndServe(":3000", nil)
}

func API_test() (ForecastResponse, error) {
	var forecastResponse ForecastResponse

	response, err := http.Get("https://api.open-meteo.com/v1/forecast?latitude=52.52&longitude=13.41&current=temperature_2m")
	if err != nil {
		return forecastResponse, fmt.Errorf("Erreur lors de la requête : %s", err)
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return forecastResponse, fmt.Errorf("Erreur lors de la lecture du corps de la réponse : %s", err)
	}

	err = json.Unmarshal(body, &forecastResponse)
	if err != nil {
		return forecastResponse, fmt.Errorf("Erreur lors du décodage JSON : %s", err)
	}
	return forecastResponse, nil
}
