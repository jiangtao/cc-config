import path from 'node:path';
import os from 'node:os';
import { cosmiconfig } from 'cosmiconfig';

export interface GitConfig {
  autoCommit: boolean;
  autoPull: boolean;
  pushRemote: string;
}

export interface BackupConfig {
  includeSettings: boolean;
  includeCommands: boolean;
  includeSkills: boolean;
}

export interface Config {
  repo: string;
  projects: string[];
  scanDirs: string[];
  lang: string;
  git: GitConfig;
  backup: BackupConfig;
  cliLanguage?: string;
}

const DEFAULTS: Config = {
  repo: '~/cc-config',
  projects: [],
  scanDirs: ['~/Places/work', '~/Places/personal', '~/work', '~/projects', '~/dev'],
  lang: 'en',
  git: {
    autoCommit: false,
    autoPull: false,
    pushRemote: 'origin',
  },
  backup: {
    includeSettings: true,
    includeCommands: true,
    includeSkills: true,
  },
};

let config: Config | null = null;

/**
 * Expand ~ to home directory
 */
export function expandPath(p: string): string {
  if (p.startsWith('~/')) {
    return path.join(os.homedir(), p.slice(2));
  }
  if (p === '~') {
    return os.homedir();
  }
  return p;
}

/**
 * Get Claude config directory
 */
export function getClaudeDir(): string {
  return path.join(os.homedir(), '.claude');
}

/**
 * Load configuration from file
 */
async function loadConfigFile(): Promise<Partial<Config>> {
  const explorer = cosmiconfig('cconf', {
    searchPlaces: [
      '.cconfrc.yaml',
      '.cconfrc.yml',
      '.cconfrc.json',
      '.cconfrc',
      'package.json',
    ],
  });

  const result = await explorer.search();
  if (result) {
    return result.config as Partial<Config>;
  }
  return {};
}

/**
 * Initialize configuration
 */
export async function initConfig(cliLanguage?: string): Promise<Config> {
  config = { ...DEFAULTS };
  const userConfig = await loadConfigFile();

  // Merge user config
  config = {
    ...config,
    ...userConfig,
    git: { ...config.git, ...userConfig.git },
    backup: { ...config.backup, ...userConfig.backup },
  };

  // Set CLI language if provided
  if (cliLanguage) {
    config.cliLanguage = cliLanguage;
  }

  return config;
}

/**
 * Get current configuration
 */
export function getConfig(): Config {
  if (!config) {
    throw new Error('Config not initialized. Call initConfig() first.');
  }
  return config;
}

/**
 * Get expanded repository path
 */
export function getRepoPath(): string {
  const cfg = getConfig();
  return expandPath(cfg.repo);
}

/**
 * Set repository path from CLI flag
 */
export function setRepo(repo: string): void {
  if (!config) {
    config = { ...DEFAULTS };
  }
  config.repo = repo;
}

/**
 * Get language for i18n
 */
export function getLanguage(): string {
  const cfg = getConfig();
  if (cfg.cliLanguage) {
    return cfg.cliLanguage;
  }
  return cfg.lang || 'en';
}

/**
 * Set language from CLI flag
 */
export function setLanguage(lang: string): void {
  if (!config) {
    config = { ...DEFAULTS };
  }
  config.cliLanguage = lang;
}
