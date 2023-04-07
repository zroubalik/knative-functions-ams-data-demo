# KubeCon Europe 2023 Demo

## Components:

The architecture is following:

1. [ams-data-app](ams-data-app) - React app that servers UI and provides HTTP endpoints to read input data and update UI accordingly. Users should specify address in Amsterdam, the app should obtain Coordinates (Latitude and Longitude). This coordinates will be used by several functions that provide some data based on the coordinates.
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

3. [scooters](scooters) - Function that accepts a cloudevent that should contain coordinates and returns list of nearest e-scooters (max 5).

```json
{
  "coordinates": {
    "latitude": 52.370216,
    "longitude": 4.895168
  },
  "scooters": [
    {
      "id": 1,
      "scooterId": "scooter-1",
      "geometry": {
        "type": "Point",
        "coordinates": [4.895168, 52.370216]
      },
      "operator": "Operator 1",
      "availabilityStatus": true,
      "name": "Scooter 1",
      "maxSpeed": 25,
      "currentLocation": "Amsterdam",
      "licensePlate": "ABC123",
      "numberOfAvailableHelmets": 2,
      "distance": 0.5
    },
    {
      "id": 2,
      "scooterId": "scooter-2",
      "geometry": {
        "type": "Point",
        "coordinates": [4.8746, 52.369553]
      },
      "operator": "Operator 2",
      "availabilityStatus": true,
      "name": "Scooter 2",
      "maxSpeed": 20,
      "currentLocation": "Amsterdam",
      "licensePlate": "XYZ789",
      "numberOfAvailableHelmets": 0,
      "distance": 1.2
    }
  ]
}
```

This JSON represents the ScooterData struct with a set of two scooters located in Amsterdam. The first scooter is available and has a current location of 0.5 kilometers from the specified coordinates. The second scooter has a current location of 1.2 kilometers from the specified coordinates.

4. (not used int he React app at the moment) [distance](distance) - Function that accepts a list of scooters as a cloudevent (a return value from Function `scooters`). And for each scooter calculates walking distance and travel time between the input coordinates and the scooter's location. Example output:

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

## How to test the React app

1. You can either run the whole app together:

```bash
cd ams-data-app
npm start
```

2. Or you can start backend and frontend separately:

Start the frontend to serve UI - it uses sockets to read from the backend

```bash
# start frontend
cd ams-data-app/frontend
npm start
```

In another termianal start the backend that listens on port :3333

```bash
# start backend
cd ams-data-app/backend
npm start
```

To display some data in the React App, you can send the example data:

```bash
curl -X POST -H "Content-Type: application/json" -d '{
  "coordinates": {
    "latitude": 52.370216,
    "longitude": 4.895168
  },
  "scooters": [
    {
      "id": 1,
      "scooterId": "scooter-1",
      "geometry": {
        "type": "Point",
        "coordinates": [4.895168, 52.370216]
      },
      "operator": "Operator 1",
      "availabilityStatus": true,
      "name": "Scooter 1",
      "maxSpeed": 25,
      "currentLocation": "Amsterdam",
      "licensePlate": "ABC123",
      "numberOfAvailableHelmets": 2,
      "distance": 0.5
    },
    {
      "id": 2,
      "scooterId": "scooter-2",
      "geometry": {
        "type": "Point",
        "coordinates": [4.8746, 52.369553]
      },
      "operator": "Operator 2",
      "availabilityStatus": true,
      "name": "Scooter 2",
      "maxSpeed": 20,
      "currentLocation": "Amsterdam",
      "licensePlate": "XYZ789",
      "numberOfAvailableHelmets": 0,
      "distance": 1.2
    }
  ]
}' http://localhost:3333/scooters

```
or
```bash
curl -k -X POST -H "Content-Type: application/json" -d '{
  "coordinates": {
    "latitude": 52.370216,
    "longitude": 4.895168
  },
  "current_weather": {
    "temperature": 20,
    "windspeed": 10,
    "weathercode": "Sunny",
    "time": "2023-04-06T14:25:00Z"
  }
}' http://localhost:3333/scooters
```
curl -k -X POST -H "Content-Type: application/json" -d '{
  "coordinates": {
    "latitude": 52.370216,
    "longitude": 4.895168
  },
  "current_weather": {
    "temperature": 20,
    "windspeed": 10,
    "weathercode": "Sunny",
    "time": "2023-04-06T14:25:00Z"
  }
}' http://ams-data-app-service:3333/weather

## How to deploy on Openshift

1. Create namespace `demo`
```bash
oc new-project demo
```

2. In [resources/app.yaml](resources/app.yaml) edit `REACT_APP_BACKEND_URI` ENV variable to match your OpenShift hostname:

```yaml
env:
  - name: REACT_APP_BACKEND_URI
    value: http://ams-data-app-api-route-demo.apps.<HOSTNAME>
```
for example:
```yaml
env:
  - name: REACT_APP_BACKEND_URI
    value: http://ams-data-app-api-route-demo.apps.zroubali.serverless.devcluster.openshift.com
```

3. Deploy resources
```bash
oc apply -f resources/
```

4. Deploy functions
```bash
func deploy -p weather
func deploy -p scooters
func deploy -p responder
```