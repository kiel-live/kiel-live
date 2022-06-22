import { createI18n } from 'vue-i18n';

import messages from '.';

function getUserLanguage(): string {
  return navigator.language.split('-')[0];
}

const i18n = createI18n({
  legacy: false,
  locale: getUserLanguage(),
  messages,
});

export default i18n;
