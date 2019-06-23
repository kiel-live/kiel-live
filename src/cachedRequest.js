const axios = require('axios');
const qs = require('querystring');

async function post(url, data) {
  const options = {
    method: 'POST',
    headers: { 'content-type': 'application/x-www-form-urlencoded' },
    data: qs.stringify(data),
    url,
  };

  try {
    const repsonse = await axios(options);
    return repsonse.data;
  } catch (e) {
    console.log('HTTP-ERROR', url, data);
    return null;
  }
}

module.exports = {
  post
}