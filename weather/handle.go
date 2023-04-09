package function

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/google/uuid"
)

type Coordinates struct {
	Latitude  string `json:"latitude"`
	Longitude string `json:"longitude"`
}

type CurrentWeather struct {
	Temperature float64 `json:"temperature"`
	Windspeed   float64 `json:"windspeed"`
	WeatherCode int     `json:"weathercode"`
	Time        string  `json:"time"`
}

type WeatherData struct {
	Coordinates    Coordinates    `json:"coordinates"`
	CurrentWeather CurrentWeather `json:"current_weather"`
}

func Handle(ctx context.Context, event cloudevents.Event) (*cloudevents.Event, error) {
	fmt.Printf("Received data %s\n", string(event.Data()))
	var coordinates Coordinates
	err := event.DataAs(&coordinates)
	if err != nil {
		return nil, fmt.Errorf("failed to get coordinates from event data: %v", err)
	}

	fmt.Printf("coordinates %+v\n", coordinates)

	url := fmt.Sprintf("https://api.open-meteo.com/v1/forecast?latitude=%s&longitude=%s&current_weather=true", coordinates.Latitude, coordinates.Longitude)
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to get weather data: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read weather data response body: %v", err)
	}

	var weatherData WeatherData
	err = json.Unmarshal(body, &weatherData)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal weather data: %v", err)
	}
	weatherData.Coordinates = coordinates

	responseEvent := cloudevents.NewEvent()
	responseEvent.SetSource("weather-provider")
	responseEvent.SetID(uuid.New().String())
	responseEvent.SetType("weather")
	responseEvent.SetData(cloudevents.ApplicationJSON, weatherData)

	fmt.Printf("%v\n", responseEvent)

	return &responseEvent, nil
}
