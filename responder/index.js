const https = require('https');
const http = require('http');
const { CloudEvent } = require('cloudevents');

const sendDataToAPI = async (url, data) => {
  const options = {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json'
    }
  };

  const client = url.startsWith('https') ? https : http;

  return new Promise((resolve, reject) => {
    const req = client.request(url, options, (res) => {
      if (res.statusCode !== 200) {
        reject(new Error(`Failed to send data to API, status code: ${res.statusCode}`));
      }

      res.setEncoding('utf8');
      let responseBody = '';
      res.on('data', (chunk) => {
        responseBody += chunk;
      });
      res.on('end', () => {
        console.log(`Data was sent to ${url}, data: ${JSON.stringify(data)}`);
        resolve();
      });
    });

    req.on('error', (error) => {
      reject(error);
    });

    req.write(JSON.stringify(data));
    req.end();
  });
};

const handle = async (context, event) => {
  console.log("Received event");
  console.log(event);

  const backendApiUri = process.env.REACT_APP_BACKEND_URI || 'http://ams-data-app-service:3333';

  switch (event.type) {
    case 'weather':
      try {
        await sendDataToAPI(`${backendApiUri}/weather`, event.data);
      } catch (error) {
        context.log.error(error);
        throw error;
      }
      break;
      case 'scooters':
        try {
          await sendDataToAPI(`${backendApiUri}/scooters`, event.data);
        } catch (error) {
          context.log.error(error);
          throw error;
        }
        break;
    default:
      console.log(`Unknown event type: ${event.type}`);
  }
};

module.exports = { handle };
