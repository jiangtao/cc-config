import { createInstance, i18n, InitOptions } from 'i18next';
import { resolve } from 'node:path';
import { readFileSync } from 'node:fs';
import { getConfig } from './config.js';

// Get locales directory - works when bundled
function getLocalesDir(): string {
  // When bundled with esbuild, __dirname will be set correctly for CJS
  // For ESM, we need to use import.meta.url
  try {
    // @ts-ignore - __dirname may not exist in strict ESM but esbuild will define it
    if (typeof __dirname !== 'undefined') {
      return resolve(__dirname, '../../locales');
    }
  } catch {}
  // Fallback for pure ESM
  return resolve('./locales');
}

let i18nInstance: i18n | null = null;

/**
 * Parse language code from environment variable
 */
function parseLangFromEnv(): string {
  const lang = process.env.LANG || process.env.LC_ALL || '';
  if (lang.toLowerCase().startsWith('zh')) {
    return 'zh';
  }
  return 'en';
}

/**
 * Get language code from various sources
 */
function getLanguageCode(): string {
  try {
    const config = getConfig();
    if (config.cliLanguage) {
      return config.cliLanguage;
    }
    return config.lang || parseLangFromEnv();
  } catch {
    return parseLangFromEnv();
  }
}

/**
 * Normalize language code
 */
function normalizeLang(lang: string): string {
  const normalized = lang.toLowerCase().replace(/[-_]/, '');
  if (normalized.startsWith('zh')) {
    return 'zh';
  }
  return 'en';
}

/**
 * Initialize i18n
 */
export async function initI18n(forceLang?: string): Promise<i18n> {
  if (i18nInstance) {
    return i18nInstance;
  }

  const langCode = normalizeLang(forceLang || getLanguageCode());

  // Read JSON files directly
  const localesDir = getLocalesDir();
  const enTranslations = JSON.parse(
    readFileSync(resolve(localesDir, 'en.json'), 'utf-8')
  );
  const zhTranslations = JSON.parse(
    readFileSync(resolve(localesDir, 'zh.json'), 'utf-8')
  );

  const options: InitOptions = {
    lng: langCode,
    fallbackLng: 'en',
    resources: {
      en: {
        translation: enTranslations,
      },
      zh: {
        translation: zhTranslations,
      },
    },
  };

  i18nInstance = createInstance(options);
  await i18nInstance.init(options);

  return i18nInstance;
}

/**
 * Translate a message
 */
export function t(key: string, params?: Record<string, unknown>): string {
  if (!i18nInstance) {
    return key;
  }
  return i18nInstance.t(key, params);
}

/**
 * Get current language
 */
export function getCurrentLanguage(): string {
  return i18nInstance?.languages[0] || 'en';
}
