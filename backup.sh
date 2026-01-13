#!/bin/bash
# Claude Code 配置自动备份脚本
# 备份配置文件到 Git 仓库（不含插件缓存）

set -e

# 颜色定义
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 脚本所在目录
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
CONFIG_DIR="$SCRIPT_DIR/config"

# Claude 配置目录
CLAUDE_DIR="$HOME/.claude"

echo -e "${BLUE}=== Claude Code 配置备份 ===${NC}"
echo ""

# 1. 备份 settings.json（移除秘钥）
echo -e "${GREEN}[1/5] 备份 settings.json（移除敏感信息）...${NC}"
if [ -f "$CLAUDE_DIR/settings.json" ]; then
    # 使用 jq 移除秘钥字段，保留其他配置
    if command -v jq &> /dev/null; then
        jq 'del(.env.ANTHROPIC_AUTH_TOKEN)' "$CLAUDE_DIR/settings.json" > "$CONFIG_DIR/settings.json"
        echo -e "${GREEN}✓ 已备份 settings.json（已移除 ANTHROPIC_AUTH_TOKEN）${NC}"
    else
        # 备用方案：使用 sed 移除秘钥行
        grep -v "ANTHROPIC_AUTH_TOKEN" "$CLAUDE_DIR/settings.json" > "$CONFIG_DIR/settings.json"
        echo -e "${GREEN}✓ 已备份 settings.json（已移除 ANTHROPIC_AUTH_TOKEN）${NC}"
        echo -e "${YELLOW}  注意：建议安装 jq 以获得更精确的 JSON 处理${NC}"
    fi
else
    echo -e "${YELLOW}⚠ 未找到 settings.json，跳过${NC}"
fi

# 2. 备份自定义命令
echo -e "${GREEN}[2/5] 备份自定义命令...${NC}"
if [ -d "$CLAUDE_DIR/commands" ] && [ "$(ls -A "$CLAUDE_DIR/commands" 2>/dev/null)" ]; then
    cp -r "$CLAUDE_DIR/commands/"* "$CONFIG_DIR/commands/" 2>/dev/null || true
    echo -e "${GREEN}✓ 已备份 $(ls "$CONFIG_DIR/commands" | wc -l) 个自定义命令${NC}"
else
    echo -e "${YELLOW}⚠ 未找到自定义命令，跳过${NC}"
fi

# 3. 备份自定义技能
echo -e "${GREEN}[3/5] 备份自定义技能...${NC}"
if [ -d "$CLAUDE_DIR/skills" ] && [ "$(ls -A "$CLAUDE_DIR/skills" 2>/dev/null)" ]; then
    cp -r "$CLAUDE_DIR/skills/"* "$CONFIG_DIR/skills/" 2>/dev/null || true
    echo -e "${GREEN}✓ 已备份 $(ls "$CONFIG_DIR/skills" 2>/dev/null | wc -l) 个自定义技能${NC}"
else
    echo -e "${YELLOW}⚠ 未找到自定义技能，跳过${NC}"
fi

# 4. 收集所有项目的 .claude 配置
echo -e "${GREEN}[4/5] 扫描并备份项目配置...${NC}"

# 常见项目目录
PROJECT_DIRS=(
    "$HOME/Places/work"
    "$HOME/Places/personal"
    "$HOME/work"
    "$HOME/projects"
    "$HOME/dev"
)

PROJECT_COUNT=0
for BASE_DIR in "${PROJECT_DIRS[@]}"; do
    if [ -d "$BASE_DIR" ]; then
        # 查找包含 .claude 目录的项目
        while IFS= read -r claude_dir; do
            PROJECT_NAME=$(basename "$(dirname "$claude_dir")")
            TARGET_DIR="$CONFIG_DIR/project-configs/$PROJECT_NAME"

            # 创建目标目录
            mkdir -p "$TARGET_DIR"

            # 复制 .claude 目录内容
            cp -r "$claude_dir"/* "$TARGET_DIR/" 2>/dev/null || true

            PROJECT_COUNT=$((PROJECT_COUNT + 1))
            echo -e "  ✓ $PROJECT_NAME"
        done < <(find "$BASE_DIR" -maxdepth 2 -type d -name ".claude" 2>/dev/null)
    fi
done

if [ $PROJECT_COUNT -eq 0 ]; then
    echo -e "${YELLOW}⚠ 未找到项目配置${NC}"
else
    echo -e "${GREEN}✓ 已备份 $PROJECT_COUNT 个项目的配置${NC}"
fi

# 5. Git 提交
echo -e "${GREEN}[5/5] Git 提交...${NC}"
cd "$SCRIPT_DIR"

# 检查是否有变化
if git diff --quiet && git diff --cached --quiet; then
    echo -e "${YELLOW}⚠ 没有变化需要提交${NC}"
else
    git add config/
    git add -u
    git commit -m "chore: backup Claude Code configuration

$(date '+%Y-%m-%d %H:%M:%S')

- Settings: $(jq -r '.enabledPlugins | keys | length' "$CONFIG_DIR/settings.json" 2>/dev/null || echo "0") plugins enabled
- Commands: $(ls "$CONFIG_DIR/commands" 2>/dev/null | wc -l)
- Skills: $(ls "$CONFIG_DIR/skills" 2>/dev/null | wc -l)
- Projects: $(ls "$CONFIG_DIR/project-configs" 2>/dev/null | wc -l)"
    echo -e "${GREEN}✓ 配置已提交到本地仓库${NC}"
    echo -e "${YELLOW}  提示：运行 'git push' 推送到远程仓库${NC}"
fi

echo ""
echo -e "${BLUE}=== 备份完成 ===${NC}"
echo ""
echo -e "${YELLOW}提示：${NC}"
echo -e "  • 插件缓存未包含在自动备份中"
echo -e "  • 如需备份插件缓存，请运行: ${GREEN}./cache-plugin.sh${NC}"
echo -e "  • 推送到远程: ${GREEN}git push${NC}"
