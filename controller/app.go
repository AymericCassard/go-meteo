package controller

import (
	"github.com/a-h/templ"
	"go-meteo/view/components"
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

func Default() templ.Component {
	return components.Hello("Go-Méteo")
}
