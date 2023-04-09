import React, { useState, useEffect } from 'react';
import io from 'socket.io-client';
import SearchIcon from "@mui/icons-material/Search";
import TextField from "@mui/material/TextField";
import Button from '@mui/material/Button';
import styles from '../src/styles/Home.module.css';
import Table from '@mui/material/Table';
import TableBody from '@mui/material/TableBody';
import TableCell, { tableCellClasses } from '@mui/material/TableCell';
import TableContainer from '@mui/material/TableContainer';
import TableHead from '@mui/material/TableHead';
import TableRow from '@mui/material/TableRow';
import Paper from '@mui/material/Paper';
import AppBar from '@mui/material/AppBar';
import Toolbar from '@mui/material/Toolbar';
import Box from '@mui/material/Box';
import ElectricScooterIcon from '@mui/icons-material/ElectricScooter';
import { styled } from '@mui/material/styles';
import { Typography } from '@mui/material';
import List from '@mui/material/List';
import ListItem from '@mui/material/ListItem';
import ListItemIcon from '@mui/material/ListItemIcon';
import ListItemText from '@mui/material/ListItemText';
import LocationOnIcon from '@mui/icons-material/LocationOn';

const socket = io('http://localhost:3333');

const StyledTableCell = styled(TableCell)(({ theme }) => ({
  [`&.${tableCellClasses.head}`]: {
    backgroundColor: theme.palette.common.black,
    color: theme.palette.common.white,
  },
  [`&.${tableCellClasses.body}`]: {
    fontSize: 14,
  },
}));

const StreetSearchField = styled(TextField)({
  '& label.Mui-focused': {
    color: 'black',
    fontSize: '20px'
  },
  '& .MuiInput-underline:after': {
    borderBottomColor: 'mediumseagreen',
  },
  '& .MuiInputBase-input': {
    color: 'black',
    fontSize: '20px'
  },
  '& .MuiOutlinedInput-input': {
    color: 'black'
  },
  '& .MuiOutlinedInput-root': {
    '& fieldset': {
      borderColor: 'black',
      // backgroundColor: 'lightyellow',
      color: '#ee0000'
    },
    '&:hover fieldset': {
      borderColor: 'mediumseagreen',
    },
    '&.Mui-focused fieldset': {
      borderColor: 'mediumseagreen',
    },
  },
});

const StyledTableRow = styled(TableRow)(({ theme }) => ({
  '&:nth-of-type(odd)': {
    backgroundColor: theme.palette.action.hover,
  },
  // hide last border
  '&:last-child td, &:last-child th': {
    border: 0,
  },
}));

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
      <div className={styles.main}>
        <div className={styles.container} style={{
                position: 'absolute',
                    top: '0',
                    left: '0',
                    width: '100%',
                    height: '100vh',
                    backgroundRepeat: 'no-repeat',
                    transition: 'background-image 2s',
                    backgroundSize: 'cover',
                    backgroundPosition: 'center',
                    opacity: '0.7',
                    backgroundImage: "url(https://images.unsplash.com/photo-1598025903112-f579d7eb72aa?ixlib=rb-4.0.3&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=2940&q=80)"
              }}>
          <div className={styles.row}>
            <div>
              <AppBar position="fixed">
                <Toolbar sx={{ justifyContent: 'space-between', backgroundColor: 'mediumseagreen'}}>
                  <Box sx={{ flex: 1 }} />
                  <ElectricScooterIcon fontSize='large' sx={{ paddingRight: '10px'}} />
                  <Typography
                    variant="h6"
                    underline="none"
                    color="inherit"
                    sx={{ fontSize: 24, flexGrow: 1 }}
                    component="div"
                  >
                  Scooty Dooby-Do
                  </Typography>
                </Toolbar>
              </AppBar>
              <Toolbar />
            </div>
            <div className="form-container">
                <div style={{ width: '60%'}}>
                  <form style={{display: 'flex', paddingBottom: '20px', margin: '100px 50px 0 50px'}}>
                    <StreetSearchField
                      id="search-bar"
                      className="text"
                      value={query}
                      onChange={handleInputChange}
                      label="Enter a street name in Amsterdam"
                      variant="outlined"
                      placeholder="Search..."
                      fullWidth/>
                    <Button
                      variant="contained"
                      onClick={handleSearch}
                      style={{ margin: '5px 0 0 10px', backgroundColor: 'mediumseagreen', height: '50px', width: '150px' }}
                      ><SearchIcon />Search</Button>
                    <Button
                      variant="contained"
                      onClick={handleClear}
                      style={{ margin: '5px 0 0 10px', backgroundColor: 'mediumseagreen', height: '50px', width: '150px' }}
                      >Clear</Button>
                  </form>
                {results.length > 0 && (
                  <List dense={true} style={{ marginLeft: '10px' }}>
                    {results.map(result => (
                      <ListItem style={{ backgroundColor: 'mediumseagreen', margin: '10px'}}>
                        <ListItemIcon>
                          <LocationOnIcon style={{ color: "#118128" }}/>
                        </ListItemIcon>
                        <ListItemText
                          primary={result.display_name}
                          onClick={() => handleSelect(result)}
                          style= {{ fontSize: 20 }}
                        />
                      </ListItem>
                    ))}
                  </List>
                )}
                {selectedResult && (
                  <TableContainer component={Paper} style={{ margin: '20px 0 0 50px', width: '70%'}}>
                    <Table sx={{ minWidth: 700 }} aria-label="customized table">
                      <TableHead>
                        <TableRow>
                          <StyledTableCell>Street Address</StyledTableCell>
                          <StyledTableCell align="right">Latitude</StyledTableCell>
                          <StyledTableCell align="right">Longitude</StyledTableCell>
                        </TableRow>
                      </TableHead>
                      <TableBody>
                          <StyledTableRow key={selectedResult.place_id}>
                            <StyledTableCell component="th" scope="row">
                              {selectedResult.display_name}
                            </StyledTableCell>
                            <StyledTableCell align="right">{selectedResult.lat}</StyledTableCell>
                            <StyledTableCell align="right">{selectedResult.lon}</StyledTableCell>
                          </StyledTableRow>
                      </TableBody>
                    </Table>
                  </TableContainer>
                )}
                {weatherData && (
                  <TableContainer component={Paper} style={{ margin: '20px 0 0 50px', width: '70%'}}>
                  <Table sx={{ minWidth: 700 }} aria-label="customized table">
                    <TableHead>
                      <TableRow>
                        <StyledTableCell>Temperature</StyledTableCell>
                        <StyledTableCell align="right">Wind Speed</StyledTableCell>
                      </TableRow>
                    </TableHead>
                    <TableBody>
                        <StyledTableRow>
                          <StyledTableCell component="th" scope="row">
                            {weatherData.current_weather.temperature}
                          </StyledTableCell>
                          <StyledTableCell align="right">{weatherData.current_weather.windspeed}</StyledTableCell>
                        </StyledTableRow>
                    </TableBody>
                  </Table>
                  </TableContainer>
                )}
                {scooterData && (
                  <TableContainer component={Paper} style={{ margin: '20px 0 0 50px', width: '70%'}}>
                  <Table sx={{ minWidth: 700 }} aria-label="customized table">
                    <TableHead>
                      <TableRow>
                        <StyledTableCell>Scooter Details</StyledTableCell>
                        <StyledTableCell align="right">Available Helmets</StyledTableCell>
                        <StyledTableCell align="right">Distance</StyledTableCell>
                      </TableRow>
                    </TableHead>
                    <TableBody>
                    {scooterData.scooters.map((scooter) => (
                        <StyledTableRow>
                          <StyledTableCell component="th" scope="row">
                            Operator: {scooter.name} | License: {scooter.licensePlate}
                          </StyledTableCell>
                          <StyledTableCell align="right">{scooter.numberOfAvailableHelmets}</StyledTableCell>
                          <StyledTableCell align="right">{scooter.distance} metres</StyledTableCell>
                        </StyledTableRow>
                      ))}
                    </TableBody>
                  </Table>
                  </TableContainer>
                )}
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}

export default App;