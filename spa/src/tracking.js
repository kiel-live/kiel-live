if (process.env.VUE_APP_MATOMO_URL && process.env.VUE_APP_MATOMO_SITE) {
  const _paq = window._paq || [];
  /* tracker methods like "setCustomDimension" should be called before "trackPageView" */
  _paq.push(['trackPageView']);
  _paq.push(['enableLinkTracking']);
  (() => {
    const u = process.env.VUE_APP_MATOMO_URL;
    _paq.push(['setTrackerUrl', `${u}matomo.php`]);
    _paq.push(['setSiteId', process.env.VUE_APP_MATOMO_SITE]);
    const d = document;
    const g = d.createElement('script');
    const s = d.getElementsByTagName('script')[0];
    g.type = 'text/javascript';
    g.async = true;
    g.defer = true;
    g.src = `${u}matomo.js`;
    s.parentNode.insertBefore(g, s);
  })();
}
