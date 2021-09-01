const config = window._env_ || {};

export default (_key, fallback) => {
  const key = `APP_${_key.toUpperCase()}`;

  if (config[key]) {
    // remove quotes
    return config[key].replace(/"/g, '');
  }

  return fallback || null;
};
