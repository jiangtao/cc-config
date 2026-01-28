'use client'

import { useState } from "react"

const oneClickInstall = "curl -fsSL https://cc-conf.vercel.app/install.sh | bash"

const installMethods = [
  {
    name: "One-Click Install",
    badge: "Recommended",
    command: oneClickInstall,
    description: "Auto-detects your platform and installs the latest version",
    highlight: true,
  },
  {
    name: "npm Global",
    badge: "",
    command: "npm install -g cc-conf",
    description: "Install using npm (requires Node.js 18+)",
    highlight: false,
  },
  {
    name: "npx (No Install)",
    badge: "Convenient",
    command: "npx cc-conf",
    description: "Run without installing (always uses latest version)",
    highlight: false,
  },
]

export function Installation() {
  const [copiedId, setCopiedId] = useState<string | null>(null)

  const handleCopy = async (command: string, id: string) => {
    await navigator.clipboard.writeText(command)
    setCopiedId(id)
    setTimeout(() => setCopiedId(null), 2000)
  }

  return (
    <section id="installation" className="py-16 sm:py-24">
      <div className="container mx-auto px-4">
        <div className="mb-8 text-center sm:mb-12">
          <h2 className="mb-4 text-2xl font-bold text-slate-900 sm:text-3xl">
            Installation
          </h2>
          <p className="text-base text-slate-600 sm:text-lg">
            Choose your preferred installation method
          </p>
        </div>

        <div className="mx-auto max-w-3xl space-y-4">
          {installMethods.map((method, index) => (
            <div
              key={index}
              className={`rounded-lg border p-4 sm:p-6 ${
                method.highlight
                  ? "border-2 border-blue-500 bg-gradient-to-r from-blue-50 to-indigo-50 shadow-lg"
                  : "border-slate-200 bg-white"
              }`}
            >
              <div className="mb-4">
                <div className="mb-3 flex flex-col gap-2 sm:flex-row sm:items-center sm:justify-between">
                  <div className="flex items-center gap-2">
                    {method.badge && (
                      <span className={`rounded-full px-3 py-1 text-xs font-semibold whitespace-nowrap ${
                        method.highlight
                          ? "bg-blue-500 text-white"
                          : "bg-slate-200 text-slate-700"
                      }`}>
                        {method.badge}
                      </span>
                    )}
                    <h3 className={`font-bold ${method.highlight ? "text-lg text-slate-900" : "text-base text-slate-900"}`}>
                      {method.name}
                    </h3>
                  </div>
                  <button
                    type="button"
                    onClick={() => handleCopy(method.command, method.name)}
                    aria-label={`Copy ${method.name} command`}
                    className={`w-full rounded px-3 py-1.5 text-xs font-medium transition-colors sm:w-auto ${
                      method.highlight
                        ? "bg-blue-500 text-white hover:bg-blue-600"
                        : "bg-slate-200 text-slate-700 hover:bg-slate-300"
                    }`}
                  >
                    {copiedId === method.name ? "✓ Copied!" : "Copy"}
                  </button>
                </div>

                {method.description && (
                  <p className={`text-sm ${method.highlight ? "text-slate-600" : "text-slate-500"}`}>
                    {method.description}
                  </p>
                )}
              </div>

              <code className={`block break-all rounded p-3 text-xs font-medium ${
                method.highlight
                  ? "bg-slate-900 text-green-400"
                  : "bg-slate-800 text-green-400"
              } sm:text-sm`}>
                {method.command}
              </code>
            </div>
          ))}
        </div>
      </div>
    </section>
  )
}
