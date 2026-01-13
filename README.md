# Claude Code 配置迁移工具

用于快速迁移 Claude Code 配置（插件、技能、命令等）到其他电脑。

## 功能

- ✅ 自动备份全局配置（settings.json）
- ✅ 备份自定义命令和技能
- ✅ 收集所有项目的 `.claude` 配置
- ✅ 手动管理插件缓存（大文件）
- ✅ Git 版本控制，安全可靠
- ✅ 敏感信息（API Token）自动排除

## 目录结构

```
cc-config/
├── .git/                    # Git 仓库
├── .gitignore               # 忽略 cache/*.tar.gz
├── backup.sh                # 自动备份脚本
├── restore.sh               # 恢复配置脚本
├── cache-plugin.sh          # 插件缓存管理
├── config/                  # Git 追踪的配置
│   ├── settings.json        # 全局配置（无 API Token）
│   ├── commands/            # 自定义命令
│   ├── skills/              # 自定义技能
│   └── project-configs/     # 项目配置
└── cache/                   # 不进 Git
    └── plugins-cache.tar.gz # 插件缓存（手动打包）
```

## 快速开始

### 首次设置（源电脑）

```bash
# 1. 克隆/初始化仓库
git clone git@github.com:jiangtao/cc-config.git ~/cc-config
cd ~/cc-config

# 2. 给脚本添加执行权限
chmod +x *.sh

# 3. 运行备份
./backup.sh

# 4. 推送到远程
git push
```

### 迁移到新电脑

```bash
# 1. 克隆仓库
git clone git@github.com:jiangtao/cc-config.git ~/cc-config
cd ~/cc-config

# 2. 给脚本添加执行权限
chmod +x *.sh

# 3. 运行恢复脚本
./restore.sh
# 按提示输入 API Token

# 4. (可选) 恢复插件缓存
./cache-plugin.sh restore
```

## 脚本说明

### backup.sh - 自动备份

备份配置文件到 Git（不含插件缓存）：

```bash
./backup.sh
```

**功能：**
- 移除 API Token 后备份 `settings.json`
- 复制自定义命令和技能
- 扫描并备份所有项目配置
- 自动 Git commit

### restore.sh - 恢复配置

从 Git 恢复配置到新电脑：

```bash
./restore.sh
```

**功能：**
- Git pull 最新配置
- 提示输入 API Token
- 恢复所有配置文件
- 可选恢复插件缓存

### cache-plugin.sh - 插件缓存管理

手动管理插件缓存（大文件）：

```bash
# 备份插件缓存
./cache-plugin.sh backup

# 恢复插件缓存
./cache-plugin.sh restore

# 清理缓存文件
./cache-plugin.sh clean
```

## 日常使用

### 修改配置后

```bash
# 自动备份并提交
./backup.sh
git push
```

### 更新插件后

```bash
# 手动打包插件缓存
./cache-plugin.sh backup

# 提交到 Git（cache 目录已被 .gitignore 排除）
git add cache/
git commit -m "update plugins cache"
git push
```

### 新电脑同步

```bash
cd ~/cc-config
git pull
./restore.sh
```

## 配置说明

### settings.json

全局配置文件，备份时自动移除：
- `ANTHROPIC_AUTH_TOKEN` - API Token（敏感信息）

其他保留的配置：
- 启用的插件
- 模型配置
- 其他环境变量

### project-configs/

自动收集项目配置，扫描路径：
- `~/Places/work/`
- `~/Places/personal/`
- `~/work/`
- `~/projects/`
- `~/dev/`

### API Token

由于安全考虑，API Token 不纳入版本控制。恢复时需手动输入：

```bash
./restore.sh
# 提示: 请输入 Anthropic API Token
```

## 注意事项

1. **API Token 安全**：秘钥需自行维护，不进 Git
2. **插件缓存**：文件较大（~1MB），手动管理
3. **项目配置**：自动扫描常见目录，手动添加其他路径需修改脚本
4. **jq 依赖**：建议安装 `jq` 用于 JSON 处理

## 依赖

- **jq**：JSON 处理工具（推荐）
  ```bash
  brew install jq  # macOS
  ```

## 许可

MIT
