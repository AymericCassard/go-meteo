package controller

import (
	// "context"
	"encoding/json"
	"fmt"
	"go-meteo/view/components"
	"io"
	"net/http"
	// "github.com/a-h/templ"
)

type VillesReponses struct {
	Results []struct {
		Name      string  `json:"name"`
		Country   string  `json:"country"`
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
	} `json:"results"`
}

type HourlyTemps struct {
	Latitude             float64 `json:"latitude"`
	Longitude            float64 `json:"longitude"`
	GenerationtimeMs     float64 `json:"generationtime_ms"`
	UtcOffsetSeconds     int     `json:"utc_offset_seconds"`
	Timezone             string  `json:"timezone"`
	TimezoneAbbreviation string  `json:"timezone_abbreviation"`
	Elevation            float64 `json:"elevation"`
	HourlyUnits          struct {
		Time          string `json:"time"`
		Temperature2M string `json:"temperature_2m"`
	} `json:"hourly_units"`
	Hourly struct {
		Time          []string  `json:"time"`
		Temperature2M []float64 `json:"temperature_2m"`
	} `json:"hourly"`
}

// return arrays of hourly values for 7 days
func getDailyTemps(temps []float64) [][]float64 {
	days := make([][]float64, len(temps)/24)
	var dailyValues []float64
	for i, temp := range temps {
		dailyValues = append(dailyValues, temp)
		if (i+1)%24 == 0 && i != 0 {
			days[i/24] = dailyValues
			dailyValues = nil
		}
	}
	return days
}

func getTemp24hAvgs(dailyValues [][]float64) []float64 {
	avgs := make([]float64, len(dailyValues))
	for i, day := range dailyValues {
		sum := 0.0
		for _, temp := range day {
			sum += temp
		}
		avgs[i] = sum / float64(len(day))
	}
	return avgs
}

func api_getMatchingVilles(sent string) (VillesReponses, error) {
	var villeResponse VillesReponses
	response, err := http.Get("https://geocoding-api.open-meteo.com/v1/search?name=" + sent + "&count=5&language=fr")
	if err != nil {
		return villeResponse, err
	}
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return villeResponse, err
	}
	err = json.Unmarshal(body, &villeResponse)
	return villeResponse, nil
}

func api_getVillesTemps(latitude, longitude string) (HourlyTemps, error) {
	var hourlyTemps HourlyTemps
	response, err := http.Get("https://api.open-meteo.com/v1/forecast?latitude=" + latitude + "&longitude=" + longitude + "&hourly=temperature_2m")
	if err != nil {
		return hourlyTemps, err
	}
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return hourlyTemps, err
	}
	err = json.Unmarshal(body, &hourlyTemps)
	return hourlyTemps, nil
}

func ReturnVilles(w http.ResponseWriter, r *http.Request) {
	villes, err := api_getMatchingVilles(r.URL.Query().Get("ville"))
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}
	buttonValues := make([]components.ButtonValues, len(villes.Results))
	for i, result := range villes.Results {
		buttonValues[i] = components.ButtonValues{
			Value:     result.Name,
			Country:   result.Country,
			Latitude:  fmt.Sprintf("%f", result.Latitude),
			Longitude: fmt.Sprintf("%f", result.Longitude),
		}
	}
	components.VilleButtonContainer(buttonValues).Render(r.Context(), w)
}

func ReturnHourlyTemps(w http.ResponseWriter, r *http.Request) {
	temps, err := api_getVillesTemps(r.URL.Query().Get("lat"), r.URL.Query().Get("lon"))
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}
	// full_comp, err = templ.ToGoHTML(context.Background(), components.DataList())
	dailyTemps := getDailyTemps(temps.Hourly.Temperature2M)
	tempsAvgs := getTemp24hAvgs(dailyTemps)
	components.DataList().Render(r.Context(), w)
	components.VilleLabel(r.URL.Query().Get("name"), r.URL.Query().Get("country")).Render(r.Context(), w)
	components.WeatherTable(dailyTemps, tempsAvgs).Render(r.Context(), w)
}
