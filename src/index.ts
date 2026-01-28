#!/usr/bin/env node
import { Command } from 'commander';
import { error } from './lib/ui.js';
import { initCommand } from './cli/init.js';
import { backupCommand } from './cli/backup.js';
import { restoreCommand } from './cli/restore.js';
import { cacheCommand } from './cli/cache.js';

const program = new Command();

program
  .name('cc-conf')
  .description('Claude Code Config - backup/restore tool')
  .version('0.1.0');

// Global language option
program.option('-l, --lang <lang>', 'Language (en/zh)', 'en');

// Init command
program.addCommand(initCommand);

// Backup command
program.addCommand(backupCommand);

// Restore command
program.addCommand(restoreCommand);

// Cache command
program.addCommand(cacheCommand);

// Parse arguments
program.parseAsync(process.argv).catch(async (err) => {
  error('Error: %s', err);
  process.exit(1);
});
