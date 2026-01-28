# cconf Landing Page

This is the official website for cconf.

---

## Quick Start

### Run locally

```bash
# Install dependencies
npm install

# Start development server
npm run dev

# Open in browser
# http://localhost:3000
```

### Build for production

```bash
# Build the site
npm run build

# Preview production build
npm start
```

---

## Deploy to Vercel

### Option 1: Vercel Dashboard (Recommended)

1. Go to [vercel.com/new](https://vercel.com/new)
2. Import your `cc-config` repository
3. Set **Root Directory** to `web`
4. Click **Deploy**

### Option 2: Vercel CLI

```bash
# Install Vercel CLI
npm install -g vercel

# Login
vercel login

# Deploy from web directory
cd web
vercel --prod
```

---

## Project Structure

```
web/
├── src/
│   ├── app/           # Next.js App Router pages
│   └── components/    # React components
├── public/            # Static files (including install.sh)
├── package.json       # Dependencies
└── tailwind.config.ts # Tailwind CSS configuration
```

---

## Tech Stack

- **Framework**: Next.js 16 (App Router)
- **Styling**: Tailwind CSS
- **Icons**: Lucide React
- **Deployment**: Vercel

---

## Need Help?

- [Next.js Documentation](https://nextjs.org/docs)
- [Tailwind CSS Documentation](https://tailwindcss.com/docs)
- [Vercel Documentation](https://vercel.com/docs)

---

## Other Languages

- **中文文档:** [README-zh.md](README-zh.md)
