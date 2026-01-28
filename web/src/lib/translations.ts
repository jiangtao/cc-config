export const translations = {
  en: {
    hero: {
      title: 'Claude Code Configuration Sync',
      subtitle: 'Backup and restore your Claude Code settings, commands, skills, and project configs across all your machines with a single command.',
      getStarted: 'Get Started',
      github: 'GitHub',
    },
    features: {
      title: 'Everything you need',
      subtitle: 'Built specifically for Claude Code power users',
      items: {
        oneCommand: {
          title: 'One-Command Backup',
          description: 'Single command backs up all your Claude Code configurations to Git.',
        },
        versionControl: {
          title: 'Version Control',
          description: 'All configs stored in Git with full history. Roll back anytime.',
        },
        secure: {
          title: 'Secure',
          description: "API tokens automatically excluded. Your secrets never leave your machine.",
        },
        crossPlatform: {
          title: 'Cross-Platform',
          description: 'Single binary works on macOS, Linux, and Windows. No dependencies.',
        },
        i18n: {
          title: 'Internationalization',
          description: 'Full support for English and 中文 (Chinese) interfaces.',
        },
        projectDiscovery: {
          title: 'Project Discovery',
          description: 'Auto-scans your common directories to find all Claude Code projects.',
        },
      },
    },
    installation: {
      title: 'Installation',
      subtitle: 'Single binary, no dependencies. One command to get started:',
      recommended: 'Recommended',
      oneClickInstall: 'One-Click Install',
      copyCommand: 'Copy Command',
      copied: 'Copied!',
      autoDetect: 'Auto-detects your platform and installs the latest version.',
      orManually: 'Or manually download for your platform:',
      buildFromSource: 'Or build from source:',
    },
    usage: {
      title: 'Usage',
      subtitle: 'Get started in under 5 minutes',
      quickStart: 'Quick Start',
      commonCommands: 'Common Commands',
    },
    footer: {
      copyright: '{year} cconf. MIT License.',
      github: 'GitHub',
      documentation: 'Documentation',
    },
  },
  zh: {
    hero: {
      title: 'Cla Code 配置同步工具',
      subtitle: '用一条命令在所有电脑间备份和恢复你的 Claude Code 设置、命令、技能和项目配置。',
      getStarted: '开始使用',
      github: 'GitHub',
    },
    features: {
      title: '你需要的一切',
      subtitle: '专为 Claude Code 高级用户打造',
      items: {
        oneCommand: {
          title: '一键备份',
          description: '单条命令即可将所有 Claude Code 配置备份到 Git。',
        },
        versionControl: {
          title: '版本控制',
          description: '所有配置存储在 Git 中，完整历史记录，随时回滚。',
        },
        secure: {
          title: '安全可靠',
          description: 'API 令牌自动排除。你的机密信息永远不会离开你的机器。',
        },
        crossPlatform: {
          title: '跨平台',
          description: '单一二进制文件，支持 macOS、Linux 和 Windows。无依赖。',
        },
        i18n: {
          title: '国际化',
          description: '完整支持英文和中文界面。',
        },
        projectDiscovery: {
          title: '项目发现',
          description: '自动扫描常用目录，查找所有 Claude Code 项目。',
        },
      },
    },
    installation: {
      title: '安装',
      subtitle: '单一二进制文件，无依赖。一条命令即可开始：',
      recommended: '推荐',
      oneClickInstall: '一键安装',
      copyCommand: '复制命令',
      copied: '已复制！',
      autoDetect: '自动检测你的平台并安装最新版本。',
      orManually: '或手动下载对应平台版本：',
      buildFromSource: '或从源码构建：',
    },
    usage: {
      title: '使用指南',
      subtitle: '5 分钟内快速上手',
      quickStart: '快速开始',
      commonCommands: '常用命令',
    },
    footer: {
      copyright: '{year} cconf. MIT 许可证。',
      github: 'GitHub',
      documentation: '文档',
    },
  },
}

export type TranslationKey = keyof typeof translations.en
