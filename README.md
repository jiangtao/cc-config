# cc-conf

> **Sync your Claude Code settings across all your computers with one command**

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![GitHub Releases](https://img.shields.io/github/v/release/jiangtao/cc-config)](https://github.com/jiangtao/cc-config/releases)

---

## What is cc-conf?

**cc-conf** is a simple tool that backs up and restores your Claude Code configurations to GitHub.

**Why do you need it?**

- ✅ You use Claude Code on multiple computers (work + personal)
- ✅ You created custom commands and don't want to lose them
- ✅ You want to share project-specific settings across machines
- ✅ You want version control for your Claude settings

---

## Quick Start

### Step 1: Install

**One-line install:**

```bash
curl -fsSL https://jiangtao.vercel.app/install.sh | bash
```

**Manual install:**

```bash
# For macOS
curl -L https://github.com/jiangtao/cc-config/releases/latest/download/cc-conf-darwin-arm64 -o cc-conf
chmod +x cc-conf
sudo mv cc-conf /usr/local/bin/

# For Linux
curl -L https://github.com/jiangtao/cc-config/releases/latest/download/cc-conf-linux-amd64 -o cc-conf
chmod +x cc-conf
sudo mv cc-conf /usr/local/bin/
```

---

### Step 2: First Time Setup

**On your main computer:**

```bash
# 1. Create a folder for your configs
mkdir -p ~/cc-config
cd ~/cc-config

# 2. Initialize git repository
git init
git remote add origin git@github.com:YOURUSERNAME/cc-config.git

# 3. Backup everything!
cc-conf backup --repo ~/cc-config

# 4. Push to GitHub
git add .
git commit -m "Initial backup"
git push -u origin main
```

**On a new computer:**

```bash
# 1. Clone your configs
git clone git@github.com:YOURUSERNAME/cc-config.git ~/cc-config

# 2. Restore everything!
cc-conf restore --repo ~/cc-config
```

That's it!

---

## What Gets Backed Up?

| Component | Description |
|-----------|-------------|
| **Settings** | Global Claude Code settings (minus API tokens) |
| **Commands** | All custom commands you created |
| **Skills** | All custom skills you installed |
| **Projects** | Project-specific Claude settings |

**Note:** API tokens are NEVER stored in Git for your security.

---

## Common Commands

### Backup

```bash
# Simple backup
cc-conf backup --repo ~/cc-config

# Backup with preview
cc-conf backup --repo ~/cc-config --dry-run

# Find ALL Claude projects automatically
cc-conf backup --all-projects
```

### Restore

```bash
# Simple restore
cc-conf restore --repo ~/cc-config

# Preview changes without applying
cc-conf restore --repo ~/cc-config --dry-run
```

---

## Optional: Configuration File

Create `~/.cc-confrc.yaml` to save your preferences:

```yaml
# Where your configs are stored
repo: ~/cc-config

# Folders to scan for Claude projects
projects:
  - ~/work
  - ~/projects
  - ~/code

# Language: en or zh
lang: en

# Auto-commit after backup
git:
  autoCommit: true
  autoPush: false
```

---

## FAQ

### Q: Is my API token safe?

**A: Yes!** API tokens are automatically removed before backup and never stored in Git.

### Q: What if I don't have a GitHub account?

**A:** You can still use cc-conf! Just use a local folder:

```bash
cc-conf backup --repo ~/Documents/my-claude-configs
```

### Q: How often should I backup?

**A:** Whenever you make changes to your Claude Code setup:

```bash
# After creating new commands or skills
cc-conf backup
```

Or set up automatic backups with cron:

```bash
# Edit crontab
crontab -e

# Add this line to backup daily at 6 PM
0 18 * * * /usr/local/bin/cc-conf backup --repo ~/cc-config
```

### Q: Can I backup just some projects?

**A:** Yes! Specify which projects to include:

```bash
cc-conf backup --projects ~/work/project1 --projects ~/work/project2
```

---

## Need Help?

- **Documentation:** [https://github.com/jiangtao/cc-config](https://github.com/jiangtao/cc-config)
- **Issues:** [https://github.com/jiangtao/cc-config/issues](https://github.com/jiangtao/cc-config/issues)

---

## Other Languages

- **中文文档:** [README-zh.md](README-zh.md)

---

## License

MIT License - see [LICENSE](LICENSE) for details.
