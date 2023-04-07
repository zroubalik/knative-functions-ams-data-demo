package function

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"sort"
	"strconv"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/google/uuid"
)

type Coordinates struct {
	Latitude  string `json:"latitude"`
	Longitude string `json:"longitude"`
}

type Scooter struct {
	ID                 int       `json:"id"`                 // ID
	ScooterID          string    `json:"scooterId"`          // Scooter ID
	Geometrie          Geometrie `json:"geometrie"`          // Geometry
	Exploitant         string    `json:"exploitant"`         // Operator
	StatusBeschikbaar  bool      `json:"statusBeschikbaar"`  // Availability status
	Naam               string    `json:"naam"`               // Name
	MaxSnelheid        int       `json:"maxSnelheid"`        // Maximum speed
	HuidigeLocatie     string    `json:"huidigeLocatie"`     // Current location
	Kenteken           string    `json:"kenteken"`           // License plate
	HelmAantalAanwezig int       `json:"helmAantalAanwezig"` // Number of available helmets
	Distance           float64   `json:"distance"`           // Distance
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

type GeometryEn struct {
	Type        string    `json:"type"`
	Coordinates []float64 `json:"coordinates"`
}

type ScooterEn struct {
	ID                       int        `json:"id"`
	ScooterID                string     `json:"scooterId"`
	Geometry                 GeometryEn `json:"geometry"`
	Operator                 string     `json:"operator"`
	AvailabilityStatus       bool       `json:"availabilityStatus"`
	Name                     string     `json:"name"`
	MaxSpeed                 int        `json:"maxSpeed"`
	CurrentLocation          string     `json:"currentLocation"`
	LicensePlate             string     `json:"licensePlate"`
	NumberOfAvailableHelmets int        `json:"numberOfAvailableHelmets"`
	Distance                 float64    `json:"distance"`
}

type ScooterDataEn struct {
	Coordinates Coordinates `json:"coordinates"`
	Scooters    []ScooterEn `json:"scooters"`
}

// Handle an event.
func Handle(ctx context.Context, event cloudevents.Event) (*cloudevents.Event, error) {
	fmt.Printf("Received data %s\n", string(event.Data()))

	var coordinates Coordinates
	err := event.DataAs(&coordinates)
	if err != nil {
		fmt.Printf("fail'\n")
		return nil, fmt.Errorf("failed to get coordinates from event data: %v", err)
	}

	lat, err := strconv.ParseFloat(coordinates.Latitude, 64)
	if err != nil {
		return nil, fmt.Errorf("error parsing latitude %v", err)
	}

	lon, err := strconv.ParseFloat(coordinates.Longitude, 64)
	if err != nil {
		return nil, fmt.Errorf("error parsing longitude %v", err)
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

	fmt.Printf("Requested scooters data for latitude %f, longitude %f\n", lat, lon)

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
			s.Distance = distance(lat, lon, s.Geometrie.Coordinates[1], s.Geometrie.Coordinates[0])
			availableScooters = append(availableScooters, s)
		}
	}

	sortByDistance(availableScooters)

	var closestScooters []Scooter
	for i := 0; i < 5 && i < len(availableScooters); i++ {
		s := availableScooters[i]
		closestScooters = append(closestScooters, s)
	}

	// translate the data
	var closestScootersEn []ScooterEn
	for _, scooter := range closestScooters {
		closestScootersEn = append(closestScootersEn, ScooterEn{
			ID:                       scooter.ID,
			ScooterID:                scooter.ScooterID,
			Geometry:                 GeometryEn(scooter.Geometrie),
			Operator:                 scooter.Exploitant,
			AvailabilityStatus:       scooter.StatusBeschikbaar,
			Name:                     scooter.Naam,
			MaxSpeed:                 scooter.MaxSnelheid,
			CurrentLocation:          scooter.HuidigeLocatie,
			LicensePlate:             scooter.Kenteken,
			NumberOfAvailableHelmets: scooter.HelmAantalAanwezig,
			Distance:                 scooter.Distance,
		})
	}

	responseEventData := ScooterDataEn{
		Scooters:    closestScootersEn,
		Coordinates: Coordinates(coordinates),
	}

	// Create a new CloudEvent with the list of scooters
	responseEvent := cloudevents.NewEvent()
	responseEvent.SetSource("scooters-lister")
	responseEvent.SetID(uuid.New().String())
	responseEvent.SetType("scooters")
	responseEvent.SetData(cloudevents.ApplicationJSON, responseEventData)

	fmt.Printf("%v\n", responseEvent)

	return &responseEvent, nil
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
