package main

// Weather structure
type Weather struct {
	City    string `json:"city"`
	Weather string `json:"weather"`
}

var weatherDb = map[string]Weather{
	"London":    {City: "London", Weather: "Rainy"},
	"Rosenheim": {City: "Rosenheim", Weather: "Clear"},
	"Munich":    {City: "Munich", Weather: "Cloudy"},
}

// GetWeather returns the weather for a given city
func GetWeather(city string) Weather {
	weather, found := weatherDb[city]
	if !found {
		weather = Weather{City: city, Weather: "Unknown"}
	}
	return weather
}
