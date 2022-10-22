import { nextTick } from 'vue';
import { createI18n } from 'vue-i18n';

function getUserLanguage(): string {
  return navigator.language.split('-')[0];
}

const userLanguage = getUserLanguage();
const i18n = createI18n({
  legacy: false,
  locale: userLanguage,
});

export const loadLocaleMessages = async (locale: string) => {
  // eslint-disable-next-line @typescript-eslint/no-unsafe-assignment
  const { default: messages } = await import(`./locales/${locale}.json`);

  i18n.global.setLocaleMessage(locale, messages);

  return nextTick();
};

export const setI18nLanguage = async (lang: string): Promise<void> => {
  if (!i18n.global.availableLocales.includes(lang)) {
    await loadLocaleMessages(lang);
  }
  i18n.global.locale.value = lang;
};

await loadLocaleMessages(userLanguage);

export default i18n;
