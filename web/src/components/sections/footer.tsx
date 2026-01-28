'use client'

import { Github } from "lucide-react"
import { useLanguage } from "@/contexts/language-context"
import { translations } from "@/lib/translations"

export function Footer() {
  const { language } = useLanguage()
  const t = translations[language]
  const currentYear = new Date().getFullYear()

  return (
    <footer className="border-t bg-slate-900 py-8 text-slate-300 sm:py-12">
      <div className="container mx-auto px-4">
        <div className="flex flex-col items-center justify-between gap-4 text-center md:flex-row md:text-left">
          <p className="text-xs sm:text-sm">
            © {currentYear} {t.footer.copyright.replace('{year}', currentYear.toString())}
          </p>
          <div className="flex items-center gap-4 sm:gap-6">
            <a
              href="https://github.com/jiangtao/cc-conf"
              className="flex items-center gap-2 hover:text-white"
              target="_blank"
              rel="noopener noreferrer"
            >
              <Github className="h-4 w-4 sm:h-5 sm:w-5" />
              <span className="text-sm">{t.footer.github}</span>
            </a>
            <a
              href={`https://github.com/jiangtao/cc-conf/blob/main/README${language === 'zh' ? '-zh' : ''}.md`}
              className="text-sm hover:text-white"
              target="_blank"
              rel="noopener noreferrer"
            >
              {t.footer.documentation}
            </a>
          </div>
        </div>
      </div>
    </footer>
  )
}
