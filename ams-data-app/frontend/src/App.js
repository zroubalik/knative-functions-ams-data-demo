import React, { useState, useEffect } from 'react';
import io from 'socket.io-client';

const socket = io('http://localhost:3333');

function App() {
  const [weatherData, setWeatherData] = useState(null);
  const [scooterData, setScooterData] = useState(null);
  const [query, setQuery] = useState('');
  const [results, setResults] = useState([]);
  const [selectedResult, setSelectedResult] = useState(null);

  useEffect(() => {
    // Listen for the 'weatherData' event
    socket.on('weatherData', (data) => {
      console.log('Received weather data:', data);
      setWeatherData(data);
    });

    // Listen for the 'scooterData' event
    socket.on('scooterData', (data) => {
      console.log('Received scooter data:', data);
      setScooterData(data);
    });

    // Clean up the event listeners
    return () => {
      socket.off('weatherData');
      socket.off('scooterData');
    };
  }, []);

  function emitCoordinatesEvent(lat, lng) {
    // Emit a 'coordinates' event to the server
    socket.emit('coordinates', { lat, lng });
  }

  function handleSearch() {
    fetch(`https://nominatim.openstreetmap.org/search?q=${query}&format=json&countrycodes=NL&city=Amsterdam`)
      .then(response => response.json())
      .then(data => {
        setResults(data);
      });
  }

  function handleClear() {
    setQuery('');
    setResults([]);
    setSelectedResult(null);
  }

  function handleSelect(result) {
    setSelectedResult(result);
    setQuery(result.display_name);
    setResults([]);
    emitCoordinatesEvent(result.lat, result.lon);
  }

  function handleInputChange(e) {
    setQuery(e.target.value);
    if (e.target.value.length > 0) {
      fetch(`https://nominatim.openstreetmap.org/search?q=${e.target.value}+Amsterdam&format=json&countrycodes=NL&city=Amsterdam`)
        .then(response => response.json())
        .then(data => {
          setResults(data);
        });
    } else {
      setResults([]);
    }
  }

  return (
    <div>
      <h1>Amsterdam City Data</h1>

      <input type="text" value={query} onChange={handleInputChange} />
      <button onClick={handleSearch}>Search</button>
      <button onClick={handleClear}>Clear</button>
      {results.length > 0 && (
        <ul>
          {results.map(result => (
            <li key={result.place_id} onClick={() => handleSelect(result)}>
              {result.display_name}
            </li>
          ))}
        </ul>
      )}
      {selectedResult && (
        <table>
          <thead>
            <tr>
              <th>Address</th>
              <th>Latitude</th>
              <th>Longitude</th>
            </tr>
          </thead>
          <tbody>
            <tr key={selectedResult.place_id}>
              <td>{selectedResult.display_name}</td>
              <td>{selectedResult.lat}</td>
              <td>{selectedResult.lon}</td>
            </tr>
          </tbody>
        </table>
      )}
      {weatherData && (
        <div>
          <h2>Weather Data</h2>
          <p>Temperature: {weatherData.current_weather.temperature}</p>
          <p>Wind speed: {weatherData.current_weather.windspeed}</p>
          <p>Weather code: {weatherData.current_weather.weathercode}</p>
        </div>
      )}
      {scooterData && (
        <div>
          <h2>Scooter Data</h2>
          <ul>
            {scooterData.scooters.map((scooter) => (
              <li key={scooter.id}>
                {scooter.name} ({scooter.licensePlate}): battery level {scooter.batteryLevel}, location ({scooter.currentLocation})
              </li>
            ))}
          </ul>
        </div>
      )}
    </div>
  );
}

export default App;