import type { Metadata } from "next"
import { Inter } from "next/font/google"
import "./globals.css"
import { LanguageProvider } from "@/contexts/language-context"
import { LanguageToggle } from "@/components/language-toggle"

const inter = Inter({ subsets: ["latin"] })

export const metadata: Metadata = {
  title: "cconf - Claude Code Configuration Sync",
  description:
    "Backup and restore your Claude Code settings, commands, skills, and project configs across all your machines.",
}

export default function RootLayout({
  children,
}: {
  children: React.ReactNode
}) {
  return (
    <html lang="en">
      <body className={inter.className}>
        <LanguageProvider>
          <div className="fixed top-4 right-4 z-50">
            <LanguageToggle />
          </div>
          {children}
        </LanguageProvider>
      </body>
    </html>
  )
}
