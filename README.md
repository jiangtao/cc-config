# ccconfig

> Claude Code configuration backup/restore tool - Sync your settings, commands, skills, and project configs across machines.

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Report Card](https://goreportcard.com/badge/github.com/jiangtao/ccconfig)](https://goreportcard.com/report/github.com/jiangtao/ccconfig)

## Why ccconfig?

When you use Claude Code across multiple computers, keeping your configurations synchronized is painful:

- Custom commands you created on your work machine
- Skills you've built or installed
- Project-specific Claude settings
- Model preferences and plugin configurations

**ccconfig solves this by backing up everything to Git and restoring it on any machine with a single command.**

## Features

- One-command backup of all Claude Code configurations
- Automatic project discovery (scans common directories)
- Git-based version control for your configs
- Plugin cache management (large files handled separately)
- Sensitive data protection (API tokens never stored in Git)
- Cross-platform single binary (macOS/Linux/Windows)
- Internationalization (English/中文)

## Quick Start

### Installation

**Option 1: Download binary (Recommended)**

```bash
# macOS (Apple Silicon)
curl -L https://github.com/jiangtao/ccconfig/releases/latest/download/ccconfig-darwin-arm64 -o ccconfig
chmod +x ccconfig
sudo mv ccconfig /usr/local/bin/

# macOS (Intel)
curl -L https://github.com/jiangtao/ccconfig/releases/latest/download/ccconfig-darwin-amd64 -o ccconfig
chmod +x ccconfig
sudo mv ccconfig /usr/local/bin/

# Linux
curl -L https://github.com/jiangtao/ccconfig/releases/latest/download/ccconfig-linux-amd64 -o ccconfig
chmod +x ccconfig
sudo mv ccconfig /usr/local/bin/
```

**Option 2: Build from source**

```bash
git clone https://github.com/jiangtao/ccconfig.git
cd ccconfig
make build
sudo make install
```

### Initial Setup

On your first computer:

```bash
# 1. Create your config repository
mkdir -p ~/cc-config
cd ~/cc-config
git init
git remote add origin git@github.com:YOURUSERNAME/cc-config.git

# 2. Backup your configurations
ccconfig backup --repo ~/cc-config

# 3. Push to GitHub
git add .
git commit -m "Initial backup"
git push -u origin main
```

### On a New Computer

```bash
# 1. Clone your config repository
git clone git@github.com:YOURUSERNAME/cc-config.git ~/cc-config

# 2. Restore everything
ccconfig restore --repo ~/cc-config
```

That's it! All your Claude Code configurations are now restored.

## Usage

### Backup Command

```bash
# Basic backup (uses repo from config file)
ccconfig backup

# Specify repo location
ccconfig backup --repo ~/cc-config

# Backup specific projects only
ccconfig backup --projects ~/work/project1 --projects ~/projects/project2

# Auto-discover ALL Claude projects on your system
ccconfig backup --all-projects

# Backup only settings (skip commands, skills, projects)
ccconfig backup --no-commands --no-skills --no-projects

# Use Chinese interface
ccconfig backup --lang zh
```

### Restore Command

```bash
# Basic restore
ccconfig restore

# Preview mode (see what would change without actually changing)
ccconfig restore --dry-run

# Skip git pull (use local state as-is)
ccconfig restore --pull=false

# Restore specific components only
ccconfig restore --no-commands --no-skills
```

### Plugin Cache Management

Plugin caches are large (~1GB) and not included in regular backups. Manage them separately:

```bash
# Backup plugin cache (creates tar.gz)
ccconfig cache backup

# Restore plugin cache
ccconfig cache restore

# Clean up cache files
ccconfig cache clean
```

## Configuration

Create `~/.ccconfig.yaml` to persist your preferences:

```yaml
# Configuration repository path
repo: ~/cc-config

# Project directories to scan
projects:
  - ~/work
  - ~/projects
  - ~/dev
  - ~/code

# Default language (en or zh)
lang: en

# Git automation
git:
  autoCommit: true   # Auto-commit after backup
  autoPush: false    # Auto-push after backup

# Backup options
backup:
  includeSettings: true    # Backup global settings.json
  includeCommands: true    # Backup custom commands
  includeSkills: true      # Backup custom skills
  includeProjects: true    # Backup project .claude configs
```

## What Gets Backed Up?

| Component | Location | Backed Up | Notes |
|-----------|----------|-----------|-------|
| Global Settings | `~/.claude/settings.json` | Yes | API token removed for security |
| Custom Commands | `~/.claude/commands/` | Yes | All custom commands |
| Custom Skills | `~/.claude/skills/` | Yes | All custom skills |
| Project Configs | `*/.claude/` | Yes | Auto-scanned from configured paths |
| Plugin Caches | `~/.claude/cache/` | Separate | Use `ccconfig cache` commands |

## Security

- **API Tokens**: Never stored in Git. You'll be prompted to enter them on restore.
- **Sensitive Data**: Settings are filtered before backup to remove secrets.
- **Private Repo**: Use a private GitHub repository for your configs.

## Advanced Usage

### GitHub Actions Auto-Restore

Set up automatic restore when configuration changes:

```yaml
# .github/workflows/restore.yml
name: Restore Claude Config
on:
  push:
    paths: ['config/**']

jobs:
  restore:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: jiangtao/ccconfig/actions/restore@main
        with:
          repo: '.'
```

### Daily Backups with Cron

```bash
# Add to crontab: crontab -e
# Backup every day at 6 PM
0 18 * * * /usr/local/bin/ccconfig backup --repo ~/cc-config
```

## Troubleshooting

**"No projects found"**
- Add your project directories to `~/.ccconfig.yaml` under `projects:`
- Or use `--all-projects` flag to scan your entire home directory

**"API Token missing after restore"**
- This is intentional! Run: `export ANTHROPIC_AUTH_TOKEN=your_token_here`
- Or add to your `~/.zshrc` or `~/.bashrc`

**Permission denied errors**
- Make sure Claude Code config directory is writable
- Check file permissions: `ls -la ~/.claude/`

## Contributing

Contributions welcome! Please see [DEVELOPMENT.md](docs/DEVELOPMENT.md) for details.

## License

MIT License - see [LICENSE](LICENSE) for details.

## Links

- [Documentation](https://ccconfig.dev)
- [GitHub](https://github.com/jiangtao/ccconfig)
- [Releases](https://github.com/jiangtao/ccconfig/releases)
