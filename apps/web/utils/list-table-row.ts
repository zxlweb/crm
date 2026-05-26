/** Shared row visuals for leads / accounts / contacts list tables */

export function initialsOf(name: string): string {
  if (!name) return '·'
  const trimmed = name.trim()
  if (!trimmed) return '·'
  const ascii = /^[\u0020-\u007E]+$/.test(trimmed)
  if (ascii) {
    const words = trimmed.split(/\s+/).filter(Boolean).slice(0, 2)
    return words.map((w) => w[0]?.toUpperCase() ?? '').join('') || trimmed[0]!.toUpperCase()
  }
  return Array.from(trimmed)[0] ?? '·'
}

function hashCode(key: string): number {
  let h = 0
  for (let i = 0; i < key.length; i += 1) {
    h = Math.trunc((h << 5) - h + (key.codePointAt(i) ?? 0))
  }
  return Math.abs(h)
}

const AVATAR_TINTS = [
  'linear-gradient(135deg, color-mix(in srgb, var(--ds-brand) 22%, transparent), color-mix(in srgb, var(--ds-brand) 8%, transparent))',
  'linear-gradient(135deg, color-mix(in srgb, var(--ds-info) 22%, transparent), color-mix(in srgb, var(--ds-info) 8%, transparent))',
  'linear-gradient(135deg, color-mix(in srgb, var(--ds-success) 22%, transparent), color-mix(in srgb, var(--ds-success) 8%, transparent))',
  'linear-gradient(135deg, color-mix(in srgb, var(--ds-warning) 24%, transparent), color-mix(in srgb, var(--ds-warning) 10%, transparent))',
  'linear-gradient(135deg, color-mix(in srgb, var(--ds-danger) 20%, transparent), color-mix(in srgb, var(--ds-danger) 8%, transparent))',
]

export function listRowAvatarStyle(id: string) {
  const idx = hashCode(id) % AVATAR_TINTS.length
  return { background: AVATAR_TINTS[idx] }
}

export function engagementBarClass(score: number): string {
  if (score >= 70) return 'bg-ds-success'
  if (score >= 40) return 'bg-ds-warning'
  return 'bg-ds-danger'
}

export function engagementTextClass(score: number): string {
  if (score >= 70) return 'text-ds-success'
  if (score >= 40) return 'text-ds-fg'
  return 'text-ds-danger'
}
