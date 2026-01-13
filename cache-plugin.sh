#!/bin/bash
# Claude Code 插件缓存备份/恢复脚本
# 手动处理插件缓存（因为文件较大且频繁变化）

set -e

# 颜色定义
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# 脚本所在目录
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
CACHE_DIR="$SCRIPT_DIR/cache"

# Claude 配置目录
CLAUDE_DIR="$HOME/.claude"
PLUGINS_CACHE_DIR="$CLAUDE_DIR/plugins/cache"
CACHE_FILE="$CACHE_DIR/plugins-cache.tar.gz"

# 显示使用说明
show_usage() {
    echo -e "${BLUE}用法: $0 [backup|restore|clean]${NC}"
    echo ""
    echo "命令:"
    echo "  backup   - 备份插件缓存到压缩包"
    echo "  restore  - 从压缩包恢复插件缓存"
    echo "  clean    - 清理缓存文件"
    echo ""
    echo "示例:"
    echo "  $0 backup    # 备份插件缓存"
    echo "  $0 restore   # 恢复插件缓存"
    echo "  $0 clean     # 清理缓存文件"
}

# 备份插件缓存
backup_cache() {
    echo -e "${BLUE}=== 备份插件缓存 ===${NC}"
    echo ""

    # 检查缓存目录是否存在
    if [ ! -d "$PLUGINS_CACHE_DIR" ]; then
        echo -e "${RED}错误：未找到插件缓存目录${NC}"
        echo "  路径: $PLUGINS_CACHE_DIR"
        exit 1
    fi

    # 创建 cache 目录
    mkdir -p "$CACHE_DIR"

    # 显示将要打包的插件
    echo -e "${GREEN}检测到的插件：${NC}"
    if [ -d "$PLUGINS_CACHE_DIR" ]; then
        for plugin_dir in "$PLUGINS_CACHE_DIR"/*; do
            if [ -d "$plugin_dir" ]; then
                plugin_name=$(basename "$plugin_dir")
                size=$(du -sh "$plugin_dir" 2>/dev/null | cut -f1)
                echo -e "  • ${plugin_name} (${size})"
            fi
        done
    fi

    # 打包
    echo ""
    echo -e "${GREEN}正在打包...${NC}"
    cd "$CLAUDE_DIR/plugins"

    # 创建临时压缩文件
    temp_file="$CACHE_FILE.tmp"
    tar -czf "$temp_file" cache/ installed_plugins.json known_marketplaces.json 2>/dev/null

    # 移动到最终位置
    mv "$temp_file" "$CACHE_FILE"

    # 显示结果
    file_size=$(du -h "$CACHE_FILE" | cut -f1)
    echo ""
    echo -e "${GREEN}✓ 备份完成${NC}"
    echo -e "  文件: $CACHE_FILE"
    echo -e "  大小: ${file_size}"
    echo ""
    echo -e "${YELLOW}提示：${NC}"
    echo -e "  • 如需提交到 Git，请运行: ${GREEN}git add cache/ && git commit -m 'update plugins cache'${NC}"
}

# 恢复插件缓存
restore_cache() {
    echo -e "${BLUE}=== 恢复插件缓存 ===${NC}"
    echo ""

    # 检查备份文件是否存在
    if [ ! -f "$CACHE_FILE" ]; then
        echo -e "${RED}错误：未找到缓存备份文件${NC}"
        echo "  路径: $CACHE_FILE"
        echo ""
        echo "请先运行: $0 backup"
        exit 1
    fi

    # 创建目标目录
    mkdir -p "$CLAUDE_DIR/plugins"

    # 显示备份信息
    file_size=$(du -h "$CACHE_FILE" | cut -f1)
    echo -e "${GREEN}备份文件信息：${NC}"
    echo -e "  文件: $CACHE_FILE"
    echo -e "  大小: ${file_size}"
    echo ""

    # 列出备份内容
    echo -e "${GREEN}备份内容：${NC}"
    tar -tzf "$CACHE_FILE" | sed 's/^/  /' | head -20
    if [ $(tar -tzf "$CACHE_FILE" | wc -l) -gt 20 ]; then
        echo "  ..."
    fi
    echo ""

    # 确认恢复
    read -p "是否恢复? (y/N): " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        echo -e "${YELLOW}已取消${NC}"
        exit 0
    fi

    # 备份现有缓存
    if [ -d "$PLUGINS_CACHE_DIR" ]; then
        backup_timestamp=$(date +%Y%m%d_%H%M%S)
        backup_dir="$CACHE_DIR/cache_backup_${backup_timestamp}"
        echo -e "${YELLOW}备份现有缓存到: $backup_dir${NC}"
        cp -r "$PLUGINS_CACHE_DIR" "$backup_dir" 2>/dev/null || true
    fi

    # 解压
    echo -e "${GREEN}正在解压...${NC}"
    tar -xzf "$CACHE_FILE" -C "$CLAUDE_DIR/plugins/"

    echo ""
    echo -e "${GREEN}✓ 恢复完成${NC}"
    echo -e "${YELLOW}提示：重启 Claude Code 使更改生效${NC}"
}

# 清理缓存
clean_cache() {
    echo -e "${BLUE}=== 清理插件缓存备份 ===${NC}"
    echo ""

    if [ ! -d "$CACHE_DIR" ]; then
        echo -e "${YELLOW}cache 目录不存在${NC}"
        exit 0
    fi

    # 列出可清理的文件
    echo -e "${GREEN}可清理的文件：${NC}"
    total_size=0
    for file in "$CACHE_DIR"/*; do
        if [ -f "$file" ]; then
            size=$(du -h "$file" | cut -f1)
            bytes=$(du -k "$file" | cut -f1)
            echo -e "  • $(basename "$file") (${size})"
            total_size=$((total_size + bytes))
        fi
    done

    if [ $total_size -eq 0 ]; then
        echo -e "${YELLOW}没有可清理的文件${NC}"
        exit 0
    fi

    total_mb=$((total_size / 1024))
    echo ""
    echo -e "总计: ${total_mb} MB"
    echo ""

    # 确认删除
    read -p "是否删除这些文件? (y/N): " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        echo -e "${YELLOW}已取消${NC}"
        exit 0
    fi

    # 删除文件
    rm -f "$CACHE_DIR"/*.tar.gz "$CACHE_DIR"/*.tmp 2>/dev/null || true

    echo -e "${GREEN}✓ 清理完成${NC}"
}

# 主逻辑
case "${1:-}" in
    backup)
        backup_cache
        ;;
    restore)
        restore_cache
        ;;
    clean)
        clean_cache
        ;;
    *)
        show_usage
        exit 1
        ;;
esac
