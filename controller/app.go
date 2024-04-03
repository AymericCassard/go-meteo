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
type Location struct {
	ID          int
	Name        string
	Latitude    float64
	Longitude   float64
	Elevation   float64
	FeatureCode string
	CountryCode string
	Admin1ID    int
	Admin2ID    int
	Admin3ID    int
	Admin4ID    int
	Timezone    string
	Population  int
	Postcodes   []string
	CountryID   int
	Country     string
	Admin1      string
	Admin2      string
	Admin3      string
	Admin4      string
}
type LocationList struct {
	Results []Location
}

func GetCurrentTempsCoordonate(long float64, lat float64) (ForecastResponse, error) {
	var forecastResponse ForecastResponse
	longStr := strconv.Itoa(int(long))
	latStr := strconv.Itoa(int(lat))
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

func GetCoordonateByCity(city string) (LocationList, error) {
	var location LocationList
	URL := "https://geocoding-api.open-meteo.com/v1/search?name=" + city + "&count=5&language=fr"
	response, err := http.Get(URL)
	if err != nil {
		return location, fmt.Errorf("Erreur lors de la requête : %s", err)
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return location, fmt.Errorf("Erreur lors de la lecture du corps de la réponse : %s", err)
	}

	err = json.Unmarshal(body, &location)
	if err != nil {
		return location, fmt.Errorf("Erreur lors du décodage JSON : %s", err)
	}
	return location, nil
}

func Default() templ.Component {
	forecast, err := GetCurrentTempsCoordonate(52.52, 13.41)
	loca, err := GetCoordonateByCity("Angers")
	if err != nil {
		fmt.Printf("Erreur : %s\n", err)
		return components.HelloError(err.Error())
	}
	temp := strconv.Itoa(int(forecast.Current.Temperature2m))
	ville := "angers"
	datalist := components.DataList([]string{"" + loca.Results[0].Name + " / " + loca.Results[0].Timezone, "" + loca.Results[1].Name + " / " + loca.Results[1].Timezone, "" + loca.Results[2].Name + " / " + loca.Results[2].Timezone})
	return components.Hello("Go-Méteo", temp, ville, datalist)
}
