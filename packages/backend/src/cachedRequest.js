const axios = require('axios');
const qs = require('querystring');
const crypto = require('crypto');

function hash(_str) {
  const str = (typeof _str !== 'string') ? JSON.stringify(_str) : _str;
  return crypto.createHash('md5').update(str).digest('hex');
}

/**
 * @credits: cache from https://medium.com/@nirlevy/debouncing-caching-axios-requests-c1b1b264e0cb
 */
class Request {
  static cache(func, _options, ...params) {
    const options = {
      ..._options,
      ...Request._defaultOptions,
    };
    const cachekey = hash(params);
    if (!Request._once[cachekey] || options.debounce === false) {
      Request._once[cachekey] = func(...params);
      return Request._once[cachekey];
    }
    setTimeout(() => (delete Request._once[cachekey]), options.ttl);
    return Request._once[cachekey];
  }

  /**
   * POST an object to the main server.
   * Requests are debounced (cached) by default within 1-second timeframe!
   * This means that multiple calls to this method with the same parameters
   * will return the result of the first call within the last 1-second.
   *
   * @method
   * @endpoint {String}          URL
   * @data     {Object}          data to POST to the server
   * @options  {Object}          debounce options: defaults are {debounce: true, ttl:1000}
   * @config   {Object}          axios config
   * @return  {Promise}          Promise
  */
  static post(endpoint, options, data, config) {
    return Request.cache(axios.post, options, endpoint, data, config);
  }

  /**
   * GET a url from the main server.
   * Requests are debounced (cached) by default within 1-second timeframe!
   * This means that multiple calls to this method with the same parameters
   * will return the result of the first call within the last 1-second.
   *
   * @method
   * @endpoint {String}          URL
   * @options  {Object}          debounce options: defaults are {debounce: true, ttl:1000}
   * @config   {Object}          axios config
   * @return  {Promise}          Promise
   */
  static get(endpoint, options, config) {
    return Request.cache(axios.get, options, endpoint, config);
  }
}

Request._defaultOptions = { debounce: true, ttl: 1000 };
Request._once = {};

async function post(url, data) {
  const options = {
    headers: { 'content-type': 'application/x-www-form-urlencoded' },
  };

  try {
    // const response = await axios(options);
    // cache every request for 3 minutes
    const response = await Request.post(url,
      {
        debounce: true,
        ttl: 10 * 60 * 1000,
      },
      qs.stringify(data), options);

    return response.data;
  } catch (e) {
    if (process.env.NODE_ENV === 'developement') {
      // eslint-disable-next-line no-console
      console.error('HTTP-ERROR', url, data, e.response.data || null);
    } else {
      // eslint-disable-next-line no-console
      console.error('HTTP-ERROR', url, data);
    }
    return null;
  }
}

module.exports = {
  post,
};
