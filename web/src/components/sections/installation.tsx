'use client'

import { useState } from "react"

const oneClickInstall = "curl -fsSL https://cconf.vercel.app/install.sh | bash"

const installSteps = [
  {
    platform: "macOS (Apple Silicon)",
    command: "curl -L https://github.com/jiangtao/cconf/releases/latest/download/cconf-darwin-arm64 -o cconf && chmod +x cconf && sudo mv cconf /usr/local/bin/",
  },
  {
    platform: "macOS (Intel)",
    command: "curl -L https://github.com/jiangtao/cconf/releases/latest/download/cconf-darwin-amd64 -o cconf && chmod +x cconf && sudo mv cconf /usr/local/bin/",
  },
  {
    platform: "Linux",
    command: "curl -L https://github.com/jiangtao/cconf/releases/latest/download/cconf-linux-amd64 -o cconf && chmod +x cconf && sudo mv cconf /usr/local/bin/",
  },
]

export function Installation() {
  const [copiedPlatform, setCopiedPlatform] = useState<string | null>(null)

  const handleCopy = async (command: string, platform: string) => {
    await navigator.clipboard.writeText(command)
    setCopiedPlatform(platform)
    setTimeout(() => setCopiedPlatform(null), 2000)
  }

  return (
    <section id="installation" className="py-16 sm:py-24">
      <div className="container mx-auto px-4">
        <div className="mb-8 text-center sm:mb-12">
          <h2 className="mb-4 text-2xl font-bold text-slate-900 sm:text-3xl">
            Installation
          </h2>
          <p className="text-base text-slate-600 sm:text-lg">
            Single binary, no dependencies. One command to get started:
          </p>
        </div>

        {/* One-Click Install - Highlighted */}
        <div className="mx-auto max-w-3xl">
          <div className="rounded-xl border-2 border-blue-500 bg-gradient-to-r from-blue-50 to-indigo-50 p-4 shadow-lg sm:p-6">
            <div className="mb-4 flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between sm:gap-2">
              <div className="flex items-center gap-2">
                <span className="rounded-full bg-blue-500 px-3 py-1 text-xs font-semibold text-white whitespace-nowrap">
                  One-Click Install
                </span>
                <h3 className="text-base font-bold text-slate-900 sm:text-lg">
                  Install cconf
                </h3>
              </div>
              <button
                type="button"
                onClick={() => handleCopy(oneClickInstall, "oneclick")}
                aria-label="Copy one-click install command"
                className="w-full rounded-lg bg-blue-500 px-4 py-2 text-sm font-medium text-white hover:bg-blue-600 transition-colors sm:w-auto"
              >
                {copiedPlatform === "oneclick" ? "✓ Copied!" : "Copy Command"}
              </button>
            </div>
            <code className="block break-all rounded-lg bg-slate-900 p-3 text-xs font-medium text-green-400 sm:p-4 sm:text-sm">
              {oneClickInstall}
            </code>
            <p className="mt-3 text-xs text-slate-600 sm:text-sm">
              Auto-detects your platform and installs the latest version.
            </p>
          </div>
        </div>
      </div>
    </section>
  )
}
