import { readFileSync, readdirSync, statSync } from 'node:fs'
import { join, relative } from 'node:path'
import { fileURLToPath } from 'node:url'
import { describe, expect, it } from 'vitest'

const FRONTEND_ROOT = fileURLToPath(new URL('..', import.meta.url))
const COMPONENTS_DIR = join(FRONTEND_ROOT, 'components')

type Layer = 'base' | 'ui' | 'common' | 'layout' | 'feature'

function resolveLayer(filePath: string): Layer | null {
  const rel = relative(COMPONENTS_DIR, filePath).replace(/\\/g, '/')
  if (rel.startsWith('base/')) return 'base'
  if (rel.startsWith('ui/')) return 'ui'
  if (rel.startsWith('common/')) return 'common'
  if (rel.startsWith('layout/')) return 'layout'
  if (rel.startsWith('feature/')) return 'feature'
  return null
}

const LAYER_RANK: Record<Layer, number> = {
  base: 0,
  ui: 1,
  common: 2,
  layout: 3,
  feature: 4,
}

const TAG_PREFIX_LAYER: Record<string, Layer> = {
  Base: 'base',
  Ui: 'ui',
  Common: 'common',
  Layout: 'layout',
  Chart: 'ui',
  Admin: 'feature',
  App: 'feature',
  Crm: 'feature',
  Leads: 'feature',
  Accounts: 'feature',
  Ai: 'feature',
  Permission: 'common',
  Login: 'feature',
  Design: 'feature',
}

function collectVueFiles(dir: string): string[] {
  const out: string[] = []
  for (const name of readdirSync(dir)) {
    const full = join(dir, name)
    if (statSync(full).isDirectory()) {
      out.push(...collectVueFiles(full))
    } else if (name.endsWith('.vue')) {
      out.push(full)
    }
  }
  return out
}

function extractUsedPrefixes(content: string): string[] {
  const prefixes = new Set<string>()
  const tagRe = /<(Base|Ui|Common|Layout|Chart|Admin|App|Crm|Leads|Accounts|Ai|Permission|Login|Design)([A-Z][A-Za-z0-9]*)/g
  let m: RegExpExecArray | null
  while ((m = tagRe.exec(content)) !== null) {
    prefixes.add(m[1])
  }
  const importRe = /from\s+['"]~\/components\/([^'"]+)['"]/g
  while ((m = importRe.exec(content)) !== null) {
    const parts = m[1].split('/')
    const head = parts[0]
    const map: Record<string, string> = {
      base: 'Base',
      ui: parts[1] === 'chart' ? 'Chart' : 'Ui',
      common: 'Permission',
      layout: 'Layout',
      feature:
        parts[1] === 'admin'
          ? 'Admin'
          : parts[1] === 'auth'
            ? 'Login'
            : parts[1] === 'app'
              ? 'App'
              : parts[1] === 'crm'
                ? 'Crm'
                : parts[1] === 'leads'
                  ? 'Leads'
                  : parts[1] === 'accounts'
                  ? 'Accounts'
                  : parts[1] === 'ai'
                  ? 'Ai'
                  : parts[1] === 'design'
                    ? 'Design'
                    : 'Feature',
      chart: 'Chart',
    }
    if (map[head]) prefixes.add(map[head])
  }
  return [...prefixes]
}

function layerForPrefix(prefix: string): Layer | undefined {
  return TAG_PREFIX_LAYER[prefix]
}

function canUseLayer(from: Layer, to: Layer): boolean {
  if (from === to) return true
  if (to === 'feature') return false
  return LAYER_RANK[to] < LAYER_RANK[from]
}

function featureDomain(relPath: string): string | null {
  if (!relPath.startsWith('feature/')) return null
  return relPath.split('/')[1] ?? null
}

/** 仅校验 apps/web/components（feature/common）；设计系统组件在 @crm/ui-kit */
describe('app component boundaries (feature/common)', () => {
  const files = collectVueFiles(COMPONENTS_DIR)

  it('maps every component file to a known layer', () => {
    const unknown = files.filter((f) => resolveLayer(f) === null)
    expect(unknown, `Unclassified: ${unknown.map((f) => relative(FRONTEND_ROOT, f)).join(', ')}`).toEqual([])
  })

  it('forbids upward dependencies (template + explicit imports)', () => {
    const violations: string[] = []

    for (const file of files) {
      const fromLayer = resolveLayer(file)!
      const rel = relative(FRONTEND_ROOT, file)
      const relPath = relative(COMPONENTS_DIR, file).replace(/\\/g, '/')
      const content = readFileSync(file, 'utf8')
      const prefixes = extractUsedPrefixes(content)

      for (const prefix of prefixes) {
        const toLayer = layerForPrefix(prefix)
        if (!toLayer) continue
        if (!canUseLayer(fromLayer, toLayer)) {
          violations.push(`${rel}: ${fromLayer} must not use <${prefix}* (${toLayer})`)
        }
      }

      if (fromLayer === 'feature') {
        const domain = featureDomain(relPath)
        if (domain === 'admin' && prefixes.includes('Login')) {
          violations.push(`${rel}: feature/admin must not use Login* (feature/auth)`)
        }
        if (domain === 'auth' && prefixes.includes('Admin')) {
          violations.push(`${rel}: feature/auth must not use Admin* (feature/admin)`)
        }
      }
    }

    expect(violations).toEqual([])
  })

  it('ui layer must not reference feature components', () => {
    const violations: string[] = []
    for (const file of files) {
      if (resolveLayer(file) !== 'ui') continue
      const content = readFileSync(file, 'utf8')
      if (/<Admin[A-Z]/.test(content) || /<Permission[A-Z]/.test(content)) {
        violations.push(relative(FRONTEND_ROOT, file))
      }
    }
    expect(violations).toEqual([])
  })
})
