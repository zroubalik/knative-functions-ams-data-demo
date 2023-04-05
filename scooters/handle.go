package function

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"sort"

	cloudevents "github.com/cloudevents/sdk-go/v2"
)

type Coordinates struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type Scooter struct {
	ID                 int       `json:"id"`
	ScooterID          string    `json:"scooterId"`
	StatusMotor        bool      `json:"statusMotor"`
	Geometrie          Geometrie `json:"geometrie"`
	Exploitant         string    `json:"exploitant"`
	StatusBeschikbaar  bool      `json:"statusBeschikbaar"`
	Naam               string    `json:"naam"`
	MaxSnelheid        int       `json:"maxSnelheid"`
	HuidigeLocatie     string    `json:"huidigeLocatie"`
	Kenteken           string    `json:"kenteken"`
	HelmVerplicht      bool      `json:"helmVerplicht"`
	HelmAantalAanwezig int       `json:"helmAantalAanwezig"`
	Distance           float64   `json:"distance"`
}

type Geometrie struct {
	Type        string    `json:"type"`
	Coordinates []float64 `json:"coordinates"`
}

type Response struct {
	Embedded struct {
		Scooters []Scooter `json:"scooters"`
	} `json:"_embedded"`
}

type ScooterData struct {
	Coordinates Coordinates `json:"coordinates"`
	Scooters    []Scooter   `json:"scooters"`
}

const (
	RAIEvent = "cloudevents.io/example/receive"
)

func Handle(ctx context.Context, event cloudevents.Event) (*cloudevents.Event, error) {
	// var coordinates Coordinates
	// err := event.DataAs(&coordinates)
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to get coordinates from event data: %v", err)
	// }

	coordinates := Coordinates{
		Latitude:  52.374,
		Longitude: 4.9,
	}

	url := "https://api.data.amsterdam.nl/v1/deelmobiliteit/scooters"

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Set("Accept-Crs", "EPSG:4326")

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	var data Response
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}

	var availableScooters []Scooter
	for _, s := range data.Embedded.Scooters {
		if s.StatusBeschikbaar {
			s.Distance = distance(coordinates.Latitude, coordinates.Longitude, s.Geometrie.Coordinates[1], s.Geometrie.Coordinates[0])
			availableScooters = append(availableScooters, s)
		}
	}

	sortByDistance(availableScooters)

	var closestScooters []Scooter
	for i := 0; i < 5 && i < len(availableScooters); i++ {
		s := availableScooters[i]
		closestScooters = append(closestScooters, s)
	}

	responseEventData := ScooterData{
		Scooters:    closestScooters,
		Coordinates: coordinates,
	}

	// Create a new CloudEvent with the list of scooters
	responseEvent := cloudevents.NewEvent(cloudevents.VersionV1)
	responseEvent.SetType("scooters.list")
	responseEvent.SetSource("scooters-lister")
	responseEvent.SetDataContentType(cloudevents.ApplicationJSON)

	responseEvent.SetData(cloudevents.ApplicationJSON, responseEventData)

	fmt.Printf("%v\n", responseEvent)

	return &event, nil
}

func distance(lat1, lon1, lat2, lon2 float64) float64 {
	const R = 6371e3 // radius of the earth in meters
	phi1 := lat1 * math.Pi / 180
	phi2 := lat2 * math.Pi / 180
	deltaPhi := (lat2 - lat1) * math.Pi / 180
	deltaLambda := (lon2 - lon1) * math.Pi / 180

	a := math.Sin(deltaPhi/2)*math.Sin(deltaPhi/2) +
		math.Cos(phi1)*math.Cos(phi2)*
			math.Sin(deltaLambda/2)*math.Sin(deltaLambda/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	d := R * c
	return d
}

func sortByDistance(scooters []Scooter) {
	sort.Slice(scooters, func(i, j int) bool {
		return scooters[i].Distance < scooters[j].Distance
	})
}
