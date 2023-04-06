import React, { useState } from 'react';

function App() {
  const [query, setQuery] = useState('');
  const [results, setResults] = useState([]);
  const [selectedResult, setSelectedResult] = useState(null);

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
    </div>
  );
}

export default App;



// use this to send the coordinates to Knative Broker, the broker URL need to be updated before this code could work
// untested version
const { CloudEvent, HTTP } = require("cloudevents-sdk");
const brokerUrl = process.env.BROKER_URI || "http://broker-ingress.knative-eventing.svc.cluster.local/demo/default"; // the URL path to the Knative broker in "demo" namespace
async function emitCoordinatesEvent(latitude, longitude) {
  // Create the CloudEvent
  const event = new CloudEvent({
    type: "coordinates",
    source: "ams-app",
    data: { latitude, longitude },
  });

  // Create the HTTP Transport binding
  const transport = new HTTP();

  // Emit the CloudEvent
  const response = await transport.emit(event, {
    method: "POST",
    brokerUrl,
    headers: {
      "Content-Type": "application/cloudevents+json",
      "Ce-Id": event.id(),
      "Ce-Source": event.source(),
      "Ce-Type": event.type(),
      "Ce-Specversion": event.specversion(),
    },
  });

  if (response.status >= 200 && response.status < 300) {
    console.log("CloudEvent emitted successfully: ", response.status, response.statusText);
  } else {
    console.error("Failed to emit CloudEvent: ", response.status, response.statusText);
  }
}