import Vue from 'vue';
import * as Sentry from '@sentry/browser';
import * as Integrations from '@sentry/integrations';

import config from '@/libs/config';

if (config('sentry_dsn')) {
  Sentry.init({
    dsn: config('sentry_dsn'),
    environment: process.env.NODE_ENV || 'development',
    release: process.env.VUE_APP_VERSION || null,
    integrations: [new Integrations.Vue({ Vue, attachProps: true })],
  });
}

export default Sentry;
