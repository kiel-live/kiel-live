import { createI18n } from 'vue-i18n';

import de from '~/locales/de.json';
import en from '~/locales/en.json';

function getUserLanguage(): string {
  return navigator.language.split('-')[0];
}

const i18n = createI18n({
  legacy: false,
  locale: getUserLanguage(),
  messages: {
    en,
    de,
  },
});

export default i18n;
