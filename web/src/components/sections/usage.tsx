export function Usage() {
  const quickStart = [
    {
      title: "1. Create your config repository",
      description: "Create a Git repository to store your Claude Code configurations",
      code: `mkdir -p ~/cc-config
cd ~/cc-config
git init
git remote add origin git@github.com:YOURUSERNAME/cc-config.git`,
    },
    {
      title: "2. Backup your configurations",
      description: "Backup all your Claude Code settings, commands, and skills",
      code: `cconf backup --repo ~/cc-config`,
    },
    {
      title: "3. Push to GitHub",
      description: "Upload your configurations to GitHub for safekeeping",
      code: `git add .
git commit -m "Initial backup"
git push -u origin main`,
    },
    {
      title: "4. On new computer, restore",
      description: "Clone and restore your configurations on any machine",
      code: `git clone git@github.com:YOURUSERNAME/cc-config.git ~/cc-config
cconf restore --repo ~/cc-config`,
    },
  ]

  const commonCommands = [
    {
      command: "cconf backup",
      description: "Backup all configurations",
    },
    {
      command: "cconf backup --all-projects",
      description: "Backup with auto-discovery of all projects",
    },
    {
      command: "cconf restore",
      description: "Restore configurations from repo",
    },
    {
      command: "cconf restore --dry-run",
      description: "Preview changes without applying",
    },
    {
      command: "cconf cache backup",
      description: "Backup plugin caches (large files)",
    },
  ]

  return (
    <section className="border-t bg-slate-50 py-16 sm:py-24">
      <div className="container mx-auto px-4">
        <div className="mb-8 text-center sm:mb-12">
          <h2 className="mb-4 text-2xl font-bold text-slate-900 sm:text-3xl">Usage</h2>
          <p className="text-base text-slate-600 sm:text-lg">
            Get started in under 5 minutes
          </p>
        </div>

        <div className="mb-12 sm:mb-16">
          <h3 className="mb-6 text-center text-xl font-semibold text-slate-900 sm:mb-8 sm:text-2xl">
            Quick Start
          </h3>
          <div className="mx-auto max-w-3xl space-y-6">
            {quickStart.map((step, index) => (
              <div
                key={step.title}
                className="relative rounded-lg border bg-white p-5 shadow-sm sm:p-6"
              >
                <div className="absolute -left-3 -top-3 flex h-8 w-8 items-center justify-center rounded-full bg-blue-500 text-sm font-bold text-white sm:h-10 sm:w-10 sm:text-base">
                  {index + 1}
                </div>
                <div className="ml-6 sm:ml-8">
                  <h4 className="mb-2 font-semibold text-slate-900 text-base sm:text-lg">
                    {step.title}
                  </h4>
                  {step.description && (
                    <p className="mb-4 text-sm text-slate-600">
                      {step.description}
                    </p>
                  )}
                  <pre className="overflow-x-auto rounded bg-slate-900 p-4 text-xs text-green-400 sm:text-sm">
                    <code>{step.code}</code>
                  </pre>
                </div>
              </div>
            ))}
          </div>
        </div>

        <div>
          <h3 className="mb-6 text-center text-xl font-semibold text-slate-900 sm:mb-8 sm:text-2xl">
            Common Commands
          </h3>
          <div className="mx-auto max-w-2xl space-y-3">
            {commonCommands.map((item) => (
              <div
                key={item.command}
                className="flex flex-col gap-1 rounded-lg border bg-white p-3 sm:flex-row sm:items-center sm:justify-between sm:p-4"
              >
                <code className="text-xs font-mono text-slate-900 sm:text-sm">
                  {item.command}
                </code>
                <span className="text-xs text-slate-600 sm:text-sm">
                  {item.description}
                </span>
              </div>
            ))}
          </div>
        </div>
      </div>
    </section>
  )
}
