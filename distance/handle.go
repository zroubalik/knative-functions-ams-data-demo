package function

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	cloudevents "github.com/cloudevents/sdk-go/v2"
)

type Coordinates struct {
	Latitude  string `json:"latitude"`
	Longitude string `json:"longitude"`
}

type Scooter struct {
	ID        int       `json:"id"`
	ScooterID string    `json:"scooterId"`
	Geometrie Geometrie `json:"geometrie"`
}

type Geometrie struct {
	Type        string    `json:"type"`
	Coordinates []float64 `json:"coordinates"`
}

type DistanceData struct {
	ScooterID   string  `json:"scooterId"`
	Distance    float64 `json:"distance"`
	WalkingTime int     `json:"walkingTime"`
	Error       string  `json:"error"`
}

type ScooterData struct {
	Coordinates Coordinates `json:"coordinates"`
	Scooters    []Scooter   `json:"scooters"`
}

func Handle(ctx context.Context, event cloudevents.Event) (*cloudevents.Event, error) {
	// Parse the CloudEvent data into ScooterData
	var data ScooterData
	if err := json.Unmarshal(event.Data(), &data); err != nil {
		return nil, fmt.Errorf("failed to unmarshal data: %s", err)
	}

	// Create an HTTP client
	client := http.Client{Timeout: time.Second * 10}

	// Initialize an empty slice of DistanceData
	var distances []DistanceData

	// Loop through the scooters and calculate the distance for each one
	for _, scooter := range data.Scooters {
		// Prepare the API request URL
		url := fmt.Sprintf("http://router.project-osrm.org/route/v1/walking/%s,%s;%f,%f?steps=false",
			data.Coordinates.Longitude, data.Coordinates.Latitude, scooter.Geometrie.Coordinates[0], scooter.Geometrie.Coordinates[1])

		// Send the API request
		resp, err := client.Get(url)
		if err != nil {
			distances = append(distances, DistanceData{ScooterID: scooter.ScooterID, Error: "failed to connect to OSRM API"})
			continue
		}

		// Parse the API response
		var result struct {
			Routes []struct {
				Duration float64 `json:"duration"`
				Distance float64 `json:"distance"`
			} `json:"routes"`
		}

		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			distances = append(distances, DistanceData{ScooterID: scooter.ScooterID, Error: "failed to parse OSRM API response"})
			continue
		}

		// Append the distance data to the slice
		distances = append(distances, DistanceData{
			ScooterID:   scooter.ScooterID,
			Distance:    result.Routes[0].Distance,      // meters
			WalkingTime: int(result.Routes[0].Duration), // seconds
		})
	}

	// Marshal the distance data into JSON
	responseData, err := json.Marshal(distances)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal response data: %s", err)
	}

	// Create a new CloudEvent with the distance data as the payload
	responseEvent := cloudevents.NewEvent(cloudevents.VersionV1)
	responseEvent.SetType("scooters.distance")
	responseEvent.SetSource("scooter-distance-calculator")
	responseEvent.SetData(cloudevents.ApplicationJSON, responseData)

	return &responseEvent, nil
}
