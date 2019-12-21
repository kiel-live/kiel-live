export default (key, fallback) => window._env_[`APP_${key.toUpperCase()}`] || fallback || null;
