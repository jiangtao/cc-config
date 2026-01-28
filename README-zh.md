# cconf

> **用一条命令在所有电脑间同步你的 Claude Code 配置**

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![GitHub Releases](https://img.shields.io/github/v/release/jiangtao/cc-config)](https://github.com/jiangtao/cc-config/releases)

---

## 什么是 cconf？

**cconf** 是一个简单的工具，可以将你的 Claude Code 配置备份到 GitHub 并随时恢复。

**为什么需要它？**

- ✅ 你在多台电脑上使用 Claude Code（工作 + 个人）
- ✅ 你创建了自定义命令，不想丢失它们
- ✅ 你想在多台机器间共享项目配置
- ✅ 你想要版本控制你的 Claude 设置

---

## 快速开始

### 第一步：安装

**一键安装：**

```bash
curl -fsSL https://jiangtao.vercel.app/install.sh | bash
```

**手动安装：**

```bash
# macOS 系统
curl -L https://github.com/jiangtao/cc-config/releases/latest/download/cconf-darwin-arm64 -o cconf
chmod +x cconf
sudo mv cconf /usr/local/bin/

# Linux 系统
curl -L https://github.com/jiangtao/cc-config/releases/latest/download/cconf-linux-amd64 -o cconf
chmod +x cconf
sudo mv cconf /usr/local/bin/
```

---

### 第二步：首次设置

**在你的主电脑上：**

```bash
# 1. 创建一个存放配置的文件夹
mkdir -p ~/cc-config
cd ~/cc-config

# 2. 初始化 git 仓库
git init
git remote add origin git@github.com:YOURUSERNAME/cc-config.git

# 3. 备份所有配置！
cconf backup --repo ~/cc-config

# 4. 推送到 GitHub
git add .
git commit -m "Initial backup"
git push -u origin main
```

**在新电脑上：**

```bash
# 1. 克隆你的配置
git clone git@github.com:YOURUSERNAME/cc-config.git ~/cc-config

# 2. 恢复所有配置！
cconf restore --repo ~/cc-config
```

就这么简单！

---

## 备份了什么？

| 组件 | 说明 |
|------|------|
| **设置** | 全局 Claude Code 设置（不包含 API 令牌） |
| **命令** | 你创建的所有自定义命令 |
| **技能** | 你安装的所有自定义技能 |
| **项目** | 项目特定的 Claude 设置 |

**注意：** 为了安全，API 令牌永远不会存储在 Git 中。

---

## 常用命令

### 备份

```bash
# 简单备份
cconf backup --repo ~/cc-config

# 备份并预览（查看将备份什么）
cconf backup --repo ~/cc-config --dry-run

# 自动查找所有 Claude 项目
cconf backup --all-projects
```

### 恢复

```bash
# 简单恢复
cconf restore --repo ~/cc-config

# 预览更改但不应用
cconf restore --repo ~/cc-config --dry-run
```

---

## 可选：配置文件

创建 `~/.cconf.yaml` 来保存你的偏好设置：

```yaml
# 你的配置存放位置
repo: ~/cc-config

# 扫描 Claude 项目的文件夹
projects:
  - ~/work
  - ~/projects
  - ~/code

# 语言：en 或 zh
lang: zh

# 备份后自动提交
git:
  autoCommit: true
  autoPush: false
```

---

## 常见问题

### 问：我的 API 令牌安全吗？

**答：是的！** API 令牌会在备份前自动移除，永远不会存储在 Git 中。

### 问：如果没有 GitHub 账号怎么办？

**答：** 你仍然可以使用 cconf！只需使用本地文件夹：

```bash
cconf backup --repo ~/Documents/my-claude-configs
```

### 问：应该多久备份一次？

**答：** 每当你修改 Claude Code 设置时：

```bash
# 创建新命令或技能后
cconf backup
```

或者设置 cron 自动备份：

```bash
# 编辑 crontab
crontab -e

# 添加这一行每天下午 6 点备份
0 18 * * * /usr/local/bin/cconf backup --repo ~/cc-config
```

### 问：我可以只备份部分项目吗？

**答：可以！** 指定要包含的项目：

```bash
cconf backup --projects ~/work/project1 --projects ~/work/project2
```

---

## 需要帮助？

- **文档：** [https://github.com/jiangtao/cc-config](https://github.com/jiangtao/cc-config)
- **问题反馈：** [https://github.com/jiangtao/cc-config/issues](https://github.com/jiangtao/cc-config/issues)

---

## 其他语言

- **English:** [README.md](README.md)

---

## 许可证

MIT License - 详见 [LICENSE](LICENSE)。
