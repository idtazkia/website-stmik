import idTranslations from '../../public/locales/id/common.json';
import enTranslations from '../../public/locales/en/common.json';

export const translations = {
  id: idTranslations,
  en: enTranslations,
};

export type Locale = keyof typeof translations;

export function getLocaleFromUrl(url: URL): Locale {
  const pathname = url.pathname;
  if (pathname.startsWith('/en')) return 'en';
  return 'id';
}

export function t(locale: Locale, key: string): string {
  const keys = key.split('.');
  let value: any = translations[locale];

  for (const k of keys) {
    value = value?.[k];
  }

  return value || key;
}

export function localizePath(locale: Locale, path: string): string {
  // Default locale (id) doesn't have prefix
  if (locale === 'id') {
    return path;
  }
  return `/${locale}${path}`;
}
