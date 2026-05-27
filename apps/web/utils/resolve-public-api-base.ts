import os from 'node:os'

/** 开发环境局域网分享：优先 .env，否则自动选本机 192.168.x / 10.x IPv4 */
export function resolvePublicApiBase(): string {
  const fromEnv = process.env.NUXT_PUBLIC_API_BASE?.trim()
  if (fromEnv) return fromEnv.replace(/\/$/, '')

  const port = process.env.NUXT_PUBLIC_API_PORT?.trim() || '8080'
  const lanIp = pickLanIPv4()
  if (lanIp) return `http://${lanIp}:${port}`
  return `http://localhost:${port}`
}

function pickLanIPv4(): string | null {
  const candidates: string[] = []
  for (const entries of Object.values(os.networkInterfaces())) {
    if (!entries) continue
    for (const net of entries) {
      const family = net.family as string | number
      const isV4 = family === 'IPv4' || family === 4
      if (!isV4 || net.internal) continue
      candidates.push(net.address)
    }
  }
  return (
    candidates.find((a) => a.startsWith('192.168.')) ??
    candidates.find((a) => a.startsWith('10.') && !a.startsWith('10.0.0.')) ??
    candidates[0] ??
    null
  )
}
