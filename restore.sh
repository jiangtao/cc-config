#!/bin/bash
# Claude Code 配置恢复脚本
# 从 Git 仓库恢复配置到新电脑

set -e

# 颜色定义
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# 脚本所在目录
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
CONFIG_DIR="$SCRIPT_DIR/config"

# Claude 配置目录
CLAUDE_DIR="$HOME/.claude"
CACHE_FILE="$SCRIPT_DIR/cache/plugins-cache.tar.gz"

echo -e "${BLUE}=== Claude Code 配置恢复 ===${NC}"
echo ""

# 检查是否在正确的目录
if [ ! -d "$CONFIG_DIR" ]; then
    echo -e "${RED}错误：未找到 config 目录${NC}"
    echo "请确保在 cc-config 项目根目录下运行此脚本"
    exit 1
fi

# 1. 更新 Git 仓库
echo -e "${GREEN}[1/6] 更新配置仓库...${NC}"
cd "$SCRIPT_DIR"
if [ -d ".git" ]; then
    git pull
    echo -e "${GREEN}✓ 配置已更新到最新版本${NC}"
else
    echo -e "${YELLOW}⚠ 未初始化 Git 仓库，跳过更新${NC}"
fi

# 2. 创建 .claude 目录
echo -e "${GREEN}[2/6] 创建 Claude 配置目录...${NC}"
mkdir -p "$CLAUDE_DIR"/{commands,skills,plugins}
echo -e "${GREEN}✓ 目录已创建${NC}"

# 3. 恢复 settings.json
echo ""
echo -e "${GREEN}[3/6] 恢复 settings.json...${NC}"
if [ -f "$CONFIG_DIR/settings.json" ]; then
    # 检查是否已有 API Token
    if [ -f "$CLAUDE_DIR/settings.json" ]; then
        EXISTING_TOKEN=$(jq -r '.env.ANTHROPIC_AUTH_TOKEN // empty' "$CLAUDE_DIR/settings.json" 2>/dev/null || echo "")
        if [ -n "$EXISTING_TOKEN" ] && [[ ! "$EXISTING_TOKEN" == "null" ]]; then
            echo -e "${YELLOW}  检测到已有的 API Token，将保留现有秘钥${NC}"
            # 合并配置，保留现有的 API Token
            if command -v jq &> /dev/null; then
                jq --arg token "$EXISTING_TOKEN" '.env.ANTHROPIC_AUTH_TOKEN = $token' "$CONFIG_DIR/settings.json" > "$CLAUDE_DIR/settings.json"
            else
                echo -e "${RED}  错误：需要安装 jq 来合并配置${NC}"
                exit 1
            fi
        else
            # 需要输入 API Token
            echo -e "${YELLOW}  需要输入 Anthropic API Token${NC}"
            read -p "  请输入 API Token (或按 Enter 跳过): " API_TOKEN
            if [ -n "$API_TOKEN" ]; then
                if command -v jq &> /dev/null; then
                    jq --arg token "$API_TOKEN" '.env.ANTHROPIC_AUTH_TOKEN = $token' "$CONFIG_DIR/settings.json" > "$CLAUDE_DIR/settings.json"
                else
                    # 备用方案：手动添加
                    cp "$CONFIG_DIR/settings.json" "$CLAUDE_DIR/settings.json"
                    echo "    \"ANTHROPIC_AUTH_TOKEN\": \"$API_TOKEN\"," >> "$CLAUDE_DIR/settings.json"
                fi
            else
                cp "$CONFIG_DIR/settings.json" "$CLAUDE_DIR/settings.json"
                echo -e "${YELLOW}  ⚠ 未设置 API Token，稍后请手动配置${NC}"
            fi
        fi
    else
        # 没有现有配置，需要输入 API Token
        read -p "  请输入 Anthropic API Token (或按 Enter 跳过): " API_TOKEN
        if [ -n "$API_TOKEN" ]; then
            if command -v jq &> /dev/null; then
                jq --arg token "$API_TOKEN" '.env.ANTHROPIC_AUTH_TOKEN = $token' "$CONFIG_DIR/settings.json" > "$CLAUDE_DIR/settings.json"
            else
                echo -e "${RED}  错误：需要安装 jq${NC}"
                exit 1
            fi
        else
            cp "$CONFIG_DIR/settings.json" "$CLAUDE_DIR/settings.json"
            echo -e "${YELLOW}  ⚠ 未设置 API Token，稍后请手动配置${NC}"
        fi
    fi
    echo -e "${GREEN}✓ settings.json 已恢复${NC}"
else
    echo -e "${YELLOW}⚠ 未找到 settings.json，跳过${NC}"
fi

# 4. 恢复自定义命令
echo ""
echo -e "${GREEN}[4/6] 恢复自定义命令...${NC}"
if [ -d "$CONFIG_DIR/commands" ] && [ "$(ls -A "$CONFIG_DIR/commands" 2>/dev/null)" ]; then
    cp -r "$CONFIG_DIR/commands/"* "$CLAUDE_DIR/commands/" 2>/dev/null || true
    chmod +x "$CLAUDE_DIR/commands"/* 2>/dev/null || true
    echo -e "${GREEN}✓ 已恢复 $(ls "$CONFIG_DIR/commands" | wc -l) 个自定义命令${NC}"
else
    echo -e "${YELLOW}⚠ 未找到自定义命令，跳过${NC}"
fi

# 5. 恢复自定义技能
echo -e "${GREEN}[5/6] 恢复自定义技能...${NC}"
if [ -d "$CONFIG_DIR/skills" ] && [ "$(ls -A "$CONFIG_DIR/skills" 2>/dev/null)" ]; then
    cp -r "$CONFIG_DIR/skills/"* "$CLAUDE_DIR/skills/" 2>/dev/null || true
    echo -e "${GREEN}✓ 已恢复 $(ls "$CONFIG_DIR/skills" 2>/dev/null | wc -l) 个自定义技能${NC}"
else
    echo -e "${YELLOW}⚠ 未找到自定义技能，跳过${NC}"
fi

# 6. 恢复插件缓存
echo ""
echo -e "${GREEN}[6/6] 恢复插件缓存...${NC}"
if [ -f "$CACHE_FILE" ]; then
    echo -e "${YELLOW}发现插件缓存文件 ($(du -h "$CACHE_FILE" | cut -f1))${NC}"
    read -p "是否恢复插件缓存? (y/N): " -n 1 -r
    echo
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        echo -e "  解压中..."
        tar -xzf "$CACHE_FILE" -C "$CLAUDE_DIR/plugins/"
        echo -e "${GREEN}✓ 插件缓存已恢复${NC}"
    else
        echo -e "${YELLOW}⚠ 已跳过插件缓存恢复${NC}"
    fi
else
    echo -e "${YELLOW}⚠ 未找到插件缓存文件，跳过${NC}"
    echo -e "  如需恢复插件缓存，请先运行: ${GREEN}./cache-plugin.sh${NC}"
fi

# 完成
echo ""
echo -e "${BLUE}=== 恢复完成 ===${NC}"
echo ""
echo -e "${YELLOW}后续步骤：${NC}"
echo -e "  1. 重启 Claude Code 使配置生效"
echo -e "  2. 如需恢复项目配置，运行: ${GREEN}./restore-project.sh <项目路径>${NC}"
echo -e "  3. 如未设置 API Token，请手动编辑 ~/.claude/settings.json"
echo ""
