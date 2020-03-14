const config = window._env_ || {};

export default (_key, fallback) => {
  const key = `APP_${_key.toUpperCase()}`;
  return config[key] || fallback || null;
};
