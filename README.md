# KubeCon Europe 2023 Demo


The architecture is following:

1. [ams-app](ams-app) - React app that servers UI and provides HTTP endpoints to read input data and update UI accordingly. Users should specify address in Amsterdam, the app should obtain Coordinates (Latitude and Longitude). This coordinates will be used by several functions that provide some data based on the coordinates.
2. [weather](weather) - Function that accepts a cloudevent that should contain coordinates and returns weather data, this is the example json:
```json
{
  "coordinates": {
    "latitude": "40.7128",
    "longitude": "-74.0060"
  },
  "current_weather": {
    "temperature": 16.2,
    "windspeed": 8.5,
    "weathercode": 200,
    "time": "2023-04-05 12:34:56"
  }
}
```
This JSON represents the WeatherData struct with a latitude of 40.7128 and a longitude of -74.0060. The current weather has a temperature of 16.2 degrees Celsius, windspeed of 8.5 m/s, weather code of 200 (thunderstorm with light rain), and the current time is 2023-04-05 12:34:56.

2. [scooters](scooters) - Function that accepts a cloudevent that should contain coordinates and returns list of nearest e-scooters (max 5).
```json
{
  "coordinates": {
    "latitude": 52.3702,
    "longitude": 4.8952
  },
  "scooters": [
    {
      "id": 1,
      "scooterId": "XYZ-123",
      "statusMotor": true,
      "geometrie": {
        "type": "Point",
        "coordinates": [
          52.3703,
          4.8949
        ]
      },
      "exploitant": "ABC Company",
      "statusBeschikbaar": true,
      "naam": "Scooter 1",
      "maxSnelheid": 25,
      "huidigeLocatie": "Amsterdam",
      "kenteken": "1234-AB",
      "helmVerplicht": true,
      "helmAantalAanwezig": 2,
      "distance": 0.2
    },
    {
      "id": 2,
      "scooterId": "ABC-456",
      "statusMotor": true,
      "geometrie": {
        "type": "Point",
        "coordinates": [
          52.3698,
          4.8955
        ]
      },
      "exploitant": "XYZ Company",
      "statusBeschikbaar": false,
      "naam": "Scooter 2",
      "maxSnelheid": 30,
      "huidigeLocatie": "Amsterdam",
      "kenteken": "5678-CD",
      "helmVerplicht": false,
      "helmAantalAanwezig": 0,
      "distance": 0.4
    }
  ]
}
```
This JSON represents the ScooterData struct with a set of two scooters located in Amsterdam. The first scooter has an ID of 1 and a scooter ID of "XYZ-123", is available and has a current location of 0.2 kilometers from the specified coordinates. The second scooter has an ID of 2 and a scooter ID of "ABC-456", is currently unavailable and has a current location of 0.4 kilometers from the specified coordinates.

3. [distance](distance) - Function that accepts a list of scooters as a cloudevent (a return value from Function `scooters`). And for each scooter calculates walking distance and travel time between the input coordinates and the scooter's location. Example output:
```json
{
  "distanceData": [
    {
      "scooterId": "XYZ-123",
      "distance": 0.5,
      "walkingTime": 10,
      "error": ""
    },
    {
      "scooterId": "ABC-456",
      "distance": 1.2,
      "walkingTime": 20,
      "error": ""
    }
  ]
}
```
This JSON represents an array of DistanceData structs with two items. The first item has a scooter ID of "XYZ-123", a distance of 0.5 kilometers, a walking time of 10 minutes, and no error. The second item has a scooter ID of "ABC-456", a distance of 1.2 kilometers, a walking time of 20 minutes, and no error.
