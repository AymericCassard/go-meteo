package controller

import (
	"encoding/json"
	"fmt"
	"go-meteo/view/components"
	"io"
	"net/http"
)

type VillesReponses struct {
	Results []struct {
		Name      string  `json:"name"`
		Country   string  `json:"country"`
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
	} `json:"results"`
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

func ReturnVilles(w http.ResponseWriter, r *http.Request) {
	villes, err := api_getMatchingVilles(r.URL.Query().Get("ville"))
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}
	for _, result := range villes.Results {
		components.DataOption(
			result.Name,
			result.Country,
			fmt.Sprintf("%f", result.Latitude),
			fmt.Sprintf("%f", result.Longitude),
		).Render(r.Context(), w)
	}
}
