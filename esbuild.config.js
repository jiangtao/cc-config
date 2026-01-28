import esbuild from 'esbuild';
import { readFile, writeFile, chmod } from 'node:fs/promises';
import { existsSync } from 'node:fs';

async function addShebang(filePath) {
  if (existsSync(filePath)) {
    let content = await readFile(filePath, 'utf-8');
    if (!content.startsWith('#!/usr/bin/env node')) {
      content = '#!/usr/bin/env node\n' + content;
      await writeFile(filePath, content);
      await chmod(filePath, 0o755);
    }
  }
}

// ESM build for npm package
esbuild.build({
  entryPoints: ['src/index.ts'],
  bundle: true,
  platform: 'node',
  target: 'node18',
  format: 'esm',
  outdir: 'dist',
  banner: {
    js: '#!/usr/bin/env node',
  },
}).catch(() => process.exit(1));

// CJS bundle for pkg binary packaging
esbuild.build({
  entryPoints: ['src/index.ts'],
  bundle: true,
  platform: 'node',
  target: 'node18',
  format: 'cjs',
  outfile: 'dist/bundle.cjs',
}).then(() => addShebang('dist/bundle.cjs')).catch(() => process.exit(1));
