import http from 'k6/http';

export const options = {
  stages: [
    { duration: '30s', target: 100 },  // Hold at 1000 RPS for 1 minute
    { duration: '30s', target: 500 }, // Ramp-up to 5000 RPS in 30 seconds
    { duration: '1m', target: 1000 },  // Hold at 5000 RPS for 1 minute
    { duration: '30s', target: 1000 }, // Ramp-up to 15000 RPS in 30 seconds
    { duration: '1m', target: 0 },     // Ramp-down to 0 RPS in 1 minutea
  ],
  thresholds: {
    http_req_duration: ['p(95)<2000'], // 95% of requests should complete within 2000ms
    'http_req_failed': ['rate<0.01'],  // Error rate should be less than 1%
  },
};

// Function to generate a random UUID-like string
function generateRandomUUID() {
  return 'xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx'.replace(/[xy]/g, function (c) {
    const r = (Math.random() * 16) | 0;
    const v = c === 'x' ? r : (r & 0x3) | 0x8;
    return v.toString(16);
  });
}

export default function () {
  const url = 'http://127.0.0.1:3000/api/v1'; // Replace with your API endpoint
  const params = {
    headers: {
      'Content-Type': 'application/json',
    },
  };

  const payload = JSON.stringify({
    url: `${generateRandomUUID()}.com`, // Generate a random URL with .com
  });

  http.post(url, payload, params);
}
