# cconf 着陆页

这是 cconf 的官方网站。

---

## 快速开始

### 本地运行

```bash
# 安装依赖
npm install

# 启动开发服务器
npm run dev

# 在浏览器中打开
# http://localhost:3000
```

### 生产构建

```bash
# 构建网站
npm run build

# 预览生产构建
npm start
```

---

## 部署到 Vercel

### 方式 1：Vercel 控制台（推荐）

1. 访问 [vercel.com/new](https://vercel.com/new)
2. 导入你的 `cc-config` 仓库
3. 将 **Root Directory** 设置为 `web`
4. 点击 **Deploy**

### 方式 2：Vercel CLI

```bash
# 安装 Vercel CLI
npm install -g vercel

# 登录
vercel login

# 从 web 目录部署
cd web
vercel --prod
```

---

## 项目结构

```
web/
├── src/
│   ├── app/           # Next.js App Router 页面
│   └── components/    # React 组件
├── public/            # 静态文件（包括 install.sh）
├── package.json       # 依赖项
└── tailwind.config.ts # Tailwind CSS 配置
```

---

## 技术栈

- **框架**: Next.js 16 (App Router)
- **样式**: Tailwind CSS
- **图标**: Lucide React
- **部署**: Vercel

---

## 需要帮助？

- [Next.js 文档](https://nextjs.org/docs)
- [Tailwind CSS 文档](https://tailwindcss.com/docs)
- [Vercel 文档](https://vercel.com/docs)

---

## 其他语言

- **English:** [README.md](README.md)
