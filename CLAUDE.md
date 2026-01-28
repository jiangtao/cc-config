# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Build & Test Commands

```bash
# Development
pnpm dev          # Run directly with tsx

# Build
pnpm build        # Compile TypeScript to dist/
pnpm typecheck    # Type check without emitting

# Testing
pnpm test         # Run Jest tests

# Code quality
pnpm lint         # Run ESLint
pnpm format       # Format with Prettier
```

## Architecture Overview

ccconfig is a Node.js CLI tool for backing up/restoring Claude Code configurations. It uses a layered architecture:

1. **CLI Layer** (`src/cli/`): Commander.js commands that parse flags and orchestrate operations
2. **Core Layer** (`src/core/`): Business logic for backup/restore operations
3. **Library Layer** (`src/lib/`): Cross-cutting utilities (config, i18n, git, ui, file)

### Key Architecture Patterns

**Configuration Merge** (`src/lib/config.ts`): The system merges configuration from multiple sources (priority order):
1. Command-line flags (`--repo`, `--lang`)
2. Config file (searches `.ccconfig.yaml`, `.ccconfig.json`, `package.json` via cosmiconfig)
3. Default values

Use `setRepo()` / `setLanguage()` to apply CLI flag overrides, and `getRepoPath()` / `getLanguage()` to read resolved values.

**Internationalization (i18n)** (`src/lib/i18n.ts`): All user-facing messages use `t(key, params)`. Language detection order:
- CLI flag `--lang` (highest priority, stored in `config.cliLanguage`)
- Config file `lang` field
- Environment `LANG`/`LC_ALL` (parses `zh_*` → `zh`)
- Defaults to English

Translation files are in `locales/en.json` and `locales/zh.json`. The i18n instance is loaded synchronously from JSON files at runtime.

**Sensitive Data Handling**: API Tokens are NEVER stored in Git:
- Backup: removes `ANTHROPIC_AUTH_TOKEN` from `settings.json` (see `src/core/backup/settings.ts:39-40`)
- Restore: preserves existing tokens, never overwrites

**Module System**: This project uses ES modules (`"type": "module"` in package.json). All imports must use `.js` extension even for TypeScript files (e.g., `import { foo } from './lib/config.js'`).

### Adding a New Command

1. Create `src/cli/<name>.ts` with a Commander.js `Command`
2. Call business logic from `src/core/` packages
3. Register in `src/index.ts` using `program.addCommand(<name>Command)`

### Adding a New Language

1. Create `locales/<lang>.json` with translations
2. Language code is auto-detected from `LANG` env (e.g., `zh_CN.UTF-8` → `zh`)

## Critical Implementation Details

**i18n.t Signature**: Call as `t('msg.id')` for simple messages or `t('msg.id', { key: val })` for template interpolation.

**Git Operations**: All git calls use `simple-git` in `src/lib/git.ts`. The system gracefully handles non-git repos (warnings, continues).

**UI Functions**: Import from `src/lib/ui.ts`:
- `title(msg)`, `success(msg)`, `warning(msg)`, `error(msg, ...args)`
- `print(color, msg)` where `color` is from the `colors` object (cyan, green, yellow, red, dim, etc.)
- `createSpinner(text)`, `prompt(question)`, `confirm(question)`

**Path Expansion**: Use `expandPath(path)` from `src/lib/config.ts` to handle `~` expansion in paths. Claude config directory is `getClaudeDir()` (returns `~/.claude`).

**Error Handling Strategy**:
- Source file missing → Warning, skip, continue
- Target permission denied → Error, exit 1
- Git operation fails → Warning, continue
- JSON parsing fails → Error, exit 1
