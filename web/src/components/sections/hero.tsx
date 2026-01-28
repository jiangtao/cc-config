'use client'

import { Button } from "@/components/ui/button"
import { Download, Github } from "lucide-react"
import { useLanguage } from "@/contexts/language-context"
import { translations } from "@/lib/translations"

export function Hero() {
  const { language } = useLanguage()
  const t = translations[language]

  return (
    <section className="container mx-auto px-4 py-16 text-center sm:py-24">
      <div className="mx-auto max-w-3xl">
        <div className="mb-6 flex justify-center">
          <code className="rounded-lg bg-slate-100 px-4 py-2 text-sm font-mono">
            cconf
          </code>
        </div>
        <h1 className="mb-6 text-3xl font-bold tracking-tight text-slate-900 sm:text-4xl md:text-5xl lg:text-6xl">
          {t.hero.title}
        </h1>
        <p className="mb-8 text-base text-slate-600 sm:text-lg md:text-xl">
          {t.hero.subtitle}
        </p>
        <div className="flex flex-col gap-3 sm:flex-row sm:justify-center sm:gap-4">
          <Button size="lg" asChild className="w-full sm:w-auto">
            <a href="#installation">
              <Download className="mr-2 h-4 w-4" />
              {t.hero.getStarted}
            </a>
          </Button>
          <Button size="lg" variant="outline" asChild className="w-full sm:w-auto">
            <a
              href="https://github.com/jiangtao/cconf"
              target="_blank"
              rel="noopener noreferrer"
            >
              <Github className="mr-2 h-4 w-4" />
              {t.hero.github}
            </a>
          </Button>
        </div>
      </div>
    </section>
  )
}
