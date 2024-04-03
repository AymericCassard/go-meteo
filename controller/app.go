package controller

import (
	"encoding/json"
	"fmt"
	"go-meteo/view/components"
	"io"
	"net/http"
	"strconv"

	"github.com/a-h/templ"
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

func GetCurrentTempsCoordonate(long float64, lat float64) (ForecastResponse, error) {
	var forecastResponse ForecastResponse
	longStr := strconv.Itoa(int(long))
	latStr := strconv.Itoa(int(lat))
	//coordonnées d'Angers
	URL := "https://api.open-meteo.com/v1/forecast?latitude=" + longStr + "&longitude=" + latStr + "&current=temperature_2m"
	response, err := http.Get(URL)
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

func Default() templ.Component {
	forecast, err := GetCurrentTempsCoordonate(52.52, 13.41)
	if err != nil {
		fmt.Printf("Erreur : %s\n", err)
		return components.HelloError(err.Error())
	}
	temp := strconv.Itoa(int(forecast.Current.Temperature2m))
	ville := "Angers"
	datalist := components.DataList([]string{"Angers", "Angouleme", "Hamburg"})
	return components.Hello("Go-Meteo", temp, ville, datalist)
}
