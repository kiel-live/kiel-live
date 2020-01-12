const path = require('path');
const flatCache = require('flat-cache');

const cache = flatCache.load('responses', path.join(__dirname, '..', '.cache'));
const noPrune = true;

module.exports = (req, res, next) => {
  const key = `__express__${req.originalUrl || req.url}`;
  const cacheContent = cache.getKey(key);

  if (cacheContent) {
    const buffer = Buffer.from(cacheContent, 'base64');
    res.end(buffer, 'binary');
  } else {
    res.sendResponse = res.send;
    res.send = (data) => {
      const buffer = Buffer.from(data);
      cache.setKey(key, buffer.toString('base64'));
      cache.save(noPrune);
      res.end(data, 'binary');
    };
    next();
  }
};
