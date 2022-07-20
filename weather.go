package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		postgresHost(), postgresPort(), postgresUser(), postgresPassword(), postgresDb())

	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		panic("Failed to connect to weather database!")
	}

	database.Table("current_weather").AutoMigrate(&CurrentWeather{})

	// some connection pool tuning
	sqlDB, _ := database.DB()
	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetConnMaxLifetime(time.Duration(5) * time.Minute)

	DB = database
}

// Weather structure
type CurrentWeather struct {
	City       string    `json:"city" gorm:"primaryKey"`
	Weather    string    `json:"weather"`
	NextUpdate time.Time `json:"-"`
}

// GetWeather returns the weather for a given city
func GetOrRetrieveWeather(city string) CurrentWeather {
	weather := CurrentWeather{}
	err := DB.Table("current_weather").First(&weather, "city = ?", city).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		// retrieve and store weather from OpenWeatherMap API
		weather = retrieveWeather(city)
		result := DB.Table("current_weather").Create(weather)
		if result.Error != nil {
			log.Fatalf("Error creating current weather %s", result.Error)
		}
	} else {
		if nextUpdate().After(weather.NextUpdate) {
			// retrieve and store weather from OpenWeatherMap API
			weather = retrieveWeather(city)
			result := DB.Table("current_weather").Save(&weather)
			if result.Error != nil {
				log.Fatalf("Error updating current weather %s", result.Error)
			}
		}
	}

	return weather
}

func nextUpdate() time.Time {
	return time.Now().AddDate(0, 0, 1)
}

type response struct {
	Name    string           `json:"name"`
	Weather []weatherDetails `json:"weather"`
	// other fields are ignored
}

type weatherDetails struct {
	Main string `json:"main"`
	// other fields are ignored
}

func retrieveWeather(city string) CurrentWeather {
	weather := CurrentWeather{City: city, Weather: "Unknown", NextUpdate: nextUpdate()}

	c := http.Client{Timeout: time.Duration(3) * time.Second}
	uri := fmt.Sprintf("%s/data/2.5/weather", weatherUri())
	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		log.Fatalf("Error creating HTTP request %s", err)
		return weather
	}

	req.Header.Add("Accept", `application/json`)
	q := req.URL.Query()
	q.Add("q", city)
	q.Add("appid", weatherAppid())
	req.URL.RawQuery = q.Encode()

	resp, err := c.Do(req)
	if err != nil {
		log.Fatalf("Error during HTTP request %s", err)
		return weather
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error reading HTTP response %s", err)
		return weather
	}

	response := response{}
	err = json.Unmarshal(body, &response)
	if err != nil {
		log.Fatalf("Error unmarshalling JSON response %s", err)
		return weather
	}

	weather.Weather = response.Weather[0].Main
	return weather
}
