#!/usr/bin/env node
// 将 docs 下所有 .md 转为同名 .html（与 .cursorrules 文档输出规范一致）
// 用法: node scripts/build-docs-html.mjs
import fs from 'fs'
import path from 'path'
import { fileURLToPath } from 'url'
import { createRequire } from 'module'

const __dirname = path.dirname(fileURLToPath(import.meta.url))
const root = path.resolve(__dirname, '..')
const docsDir = path.join(root, 'docs')
const require = createRequire(path.join(__dirname, 'doc-tools/package.json'))

let marked
try {
  marked = require('marked')
} catch {
  console.error('请先执行: cd scripts/doc-tools && npm install')
  process.exit(1)
}

marked.setOptions({ gfm: true, breaks: false })

const CSS = `
:root {
  --bg: #f8fafc; --fg: #0f172a; --muted: #64748b;
  --card: #fff; --border: #e2e8f0; --accent: #2563eb;
  --code-bg: #f1f5f9;
}
[data-theme="dark"] {
  --bg: #0f172a; --fg: #e2e8f0; --muted: #94a3b8;
  --card: #1e293b; --border: #334155; --accent: #60a5fa;
  --code-bg: #1e293b;
}
* { box-sizing: border-box; }
body {
  margin: 0; font-family: ui-sans-serif, system-ui, -apple-system, "Segoe UI", sans-serif;
  background: var(--bg); color: var(--fg); line-height: 1.65;
}
.layout { display: flex; min-height: 100vh; }
nav {
  width: 260px; flex-shrink: 0; padding: 1.25rem; border-right: 1px solid var(--border);
  background: var(--card); position: sticky; top: 0; height: 100vh; overflow-y: auto;
}
nav h2 { font-size: 0.75rem; text-transform: uppercase; letter-spacing: .05em; color: var(--muted); margin: 1rem 0 .5rem; }
nav a { display: block; font-size: 0.875rem; color: var(--fg); text-decoration: none; padding: .25rem 0; }
nav a:hover { color: var(--accent); }
main { flex: 1; max-width: 900px; padding: 2rem 2.5rem 4rem; }
.toolbar { display: flex; gap: .75rem; margin-bottom: 1.5rem; align-items: center; }
.toolbar button, .toolbar a {
  font-size: 0.8125rem; padding: .4rem .75rem; border-radius: 6px;
  border: 1px solid var(--border); background: var(--card); color: var(--fg); cursor: pointer; text-decoration: none;
}
article h1 { font-size: 1.75rem; border-bottom: 1px solid var(--border); padding-bottom: .5rem; }
article h2 { font-size: 1.35rem; margin-top: 2rem; }
article h3 { font-size: 1.1rem; }
article table { width: 100%; border-collapse: collapse; margin: 1rem 0; font-size: 0.9rem; }
article th, article td { border: 1px solid var(--border); padding: .5rem .75rem; text-align: left; }
article th { background: var(--code-bg); }
article code { background: var(--code-bg); padding: .15em .35em; border-radius: 4px; font-size: 0.88em; }
article pre { background: var(--code-bg); padding: 1rem; border-radius: 8px; overflow-x: auto; }
article pre code { background: none; padding: 0; }
article blockquote { border-left: 4px solid var(--accent); margin: 1rem 0; padding-left: 1rem; color: var(--muted); }
article a { color: var(--accent); }
.footer { margin-top: 3rem; font-size: 0.75rem; color: var(--muted); }
`

function walkMd(dir, list = []) {
  for (const name of fs.readdirSync(dir)) {
    const p = path.join(dir, name)
    if (fs.statSync(p).isDirectory() && !name.startsWith('_')) walkMd(p, list)
    else if (name.endsWith('.md')) list.push(p)
  }
  return list.sort()
}

/** 删除已无对应 .md 的孤儿 .html（避免删 md 后 HTML 仍留在仓库） */
function pruneOrphanHtml(dir) {
  let removed = 0
  for (const name of fs.readdirSync(dir)) {
    const p = path.join(dir, name)
    if (fs.statSync(p).isDirectory() && !name.startsWith('_')) {
      removed += pruneOrphanHtml(p)
    } else if (name.endsWith('.html') && name !== 'index.html') {
      const mdPath = p.replace(/\.html$/, '.md')
      if (!fs.existsSync(mdPath)) {
        fs.unlinkSync(p)
        console.log('🗑', path.relative(root, p), '(无对应 .md)')
        removed++
      }
    }
  }
  return removed
}

function buildNav(allMd, currentRel) {
  const byDir = {}
  for (const f of allMd) {
    const rel = path.relative(docsDir, f)
    const dir = path.dirname(rel) || '.'
    if (!byDir[dir]) byDir[dir] = []
    byDir[dir].push(rel)
  }
  let html = '<p><a href="./index.html"><strong>📚 文档首页</strong></a></p>'
  for (const dir of Object.keys(byDir).sort()) {
    html += `<h2>${dir === '.' ? '根目录' : dir}</h2>`
    for (const rel of byDir[dir]) {
      const href = rel.replace(/\.md$/, '.html')
      const cur = rel === currentRel ? ' style="color:var(--accent);font-weight:600"' : ''
      html += `<a href="${path.posix.join('../'.repeat(currentRel.split('/').length - 1), href)}"${cur}>${path.basename(rel)}</a>`
    }
  }
  return html
}

function relNavPath(fromDir, toRelHtml) {
  const depth = fromDir === '.' ? 0 : fromDir.split('/').length
  return (depth ? '../'.repeat(depth) : './') + toRelHtml
}

function toPosixPath(p) {
  return p.replaceAll('\\', '/')
}

function buildSidebar(allMd, currentRel) {
  const currentRelNorm = toPosixPath(currentRel)
  const byDir = {}
  for (const f of allMd) {
    const relNorm = toPosixPath(path.relative(docsDir, f))
    const dirNorm = path.posix.dirname(relNorm) || '.'
    if (!byDir[dirNorm]) byDir[dirNorm] = []
    byDir[dirNorm].push(relNorm)
  }
  let html = `<p><a href="/index.html"><strong>📚 文档首页</strong></a></p>`
  for (const dir of Object.keys(byDir).sort()) {
    html += `<h2>${dir === '.' ? '根目录' : dir}</h2>`
    for (const rel of byDir[dir]) {
      const href = '/' + rel.replace(/\.md$/, '.html')
      const cur = rel === currentRelNorm ? ' style="color:var(--accent);font-weight:600"' : ''
      html += `<a href="${href}"${cur}>${path.posix.basename(rel, '.md')}</a>`
    }
  }
  return html
}

function wrapPage(title, bodyHtml, sidebar, relFromDocs) {
  return `<!DOCTYPE html>
<html lang="zh-CN">
<head>
  <meta charset="UTF-8" />
  <meta name="viewport" content="width=device-width, initial-scale=1" />
  <title>${title} · CRM Docs</title>
  <style>${CSS}</style>
  <script src="https://cdn.jsdelivr.net/npm/mermaid@11/dist/mermaid.min.js"></script>
  <link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/highlight.js@11/styles/github.min.css" id="hljs-light" />
  <link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/highlight.js@11/styles/github-dark.min.css" id="hljs-dark" disabled />
  <script src="https://cdn.jsdelivr.net/gh/highlightjs/cdn-release@11/build/highlight.min.js"></script>
</head>
<body>
  <div class="layout">
    <nav>${sidebar}</nav>
    <main>
      <div class="toolbar">
        <button type="button" id="theme-toggle">🌓 切换主题</button>
        <a href="${relFromDocs.replace(/\.html$/, '.md')}" target="_blank">查看 Markdown 源码</a>
      </div>
      <article class="markdown-body">${bodyHtml}</article>
      <p class="footer">由 scripts/build-docs-html.mjs 生成 · 源文件 ${relFromDocs}</p>
    </main>
  </div>
  <script>
    mermaid.initialize({ startOnLoad: true, theme: document.documentElement.dataset.theme === 'dark' ? 'dark' : 'default' });
    document.querySelectorAll('pre code').forEach(el => { try { hljs.highlightElement(el); } catch(e) {} });
    const key = 'crm-docs-theme';
    const saved = localStorage.getItem(key);
    if (saved === 'dark') document.documentElement.dataset.theme = 'dark';
    document.getElementById('theme-toggle').onclick = () => {
      const d = document.documentElement.dataset.theme === 'dark' ? '' : 'dark';
      document.documentElement.dataset.theme = d;
      localStorage.setItem(key, d || 'light');
      document.getElementById('hljs-dark').disabled = d !== 'dark';
      document.getElementById('hljs-light').disabled = d === 'dark';
    };
    if (document.documentElement.dataset.theme === 'dark') {
      document.getElementById('hljs-dark').disabled = false;
      document.getElementById('hljs-light').disabled = true;
    }
  </script>
</body>
</html>`
}

const allMd = walkMd(docsDir)
const pruned = pruneOrphanHtml(docsDir)
let count = 0

for (const mdPath of allMd) {
  const rel = path.relative(docsDir, mdPath)
  const htmlPath = mdPath.replace(/\.md$/, '.html')
  const raw = fs.readFileSync(mdPath, 'utf8')
  const body = marked.parse(raw)
  const title = raw.match(/^#\s+(.+)/m)?.[1] || path.basename(rel, '.md')
  const depth = path.dirname(rel) === '.' ? 0 : path.dirname(rel).split('/').length
  const prefix = depth ? '../'.repeat(depth) : './'
  const sidebar = buildSidebar(allMd, rel)
  fs.writeFileSync(htmlPath, wrapPage(title, body, sidebar, prefix + rel))
  count++
  console.log('✓', rel, '→', path.relative(root, htmlPath))
}

// docs/index.html 入口
const indexMd = path.join(docsDir, 'README.md')
const indexBody = marked.parse(fs.readFileSync(indexMd, 'utf8'))
const indexSidebar = buildSidebar(allMd, 'README.md')
fs.writeFileSync(path.join(docsDir, 'index.html'), wrapPage('CRM 项目文档', indexBody, indexSidebar, './'))

console.log(`\n✅ 已生成 ${count} 个 HTML + docs/index.html` + (pruned ? `，清理孤儿 HTML ${pruned} 个` : ''))
