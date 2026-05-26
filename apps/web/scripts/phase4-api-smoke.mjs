const base = process.env.API_BASE || 'http://localhost:8080'
const tenant = 'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa'

async function main() {
  const loginRes = await fetch(`${base}/api/auth/login`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ email: 'admin@demo.com', password: 'password123' }),
  })
  const login = await loginRes.json()
  const token = login.data?.access_token
  if (!token) throw new Error('login failed: ' + JSON.stringify(login))

  const tenantHeaders = {
    Authorization: `Bearer ${token}`,
    'X-Tenant-ID': tenant,
  }

  const checks = [
    ['GET /api/settings/tenant', `${base}/api/settings/tenant`, tenantHeaders],
    ['GET /api/settings/custom-fields', `${base}/api/settings/custom-fields?entity_type=lead`, tenantHeaders],
    ['GET /api/audit/stats/by-action', `${base}/api/audit/stats/by-action?from=2026-05-01&to=2026-05-26`, tenantHeaders],
    ['GET /api/super-admin/stats/tenant-health', `${base}/api/super-admin/stats/tenant-health`, { Authorization: `Bearer ${token}` }],
  ]

  for (const [name, url, headers] of checks) {
    const res = await fetch(url, { headers })
    const body = await res.json()
    console.log(name, res.status, body.code, body.message || '')
    if (!res.ok || body.code !== 200) {
      console.error(JSON.stringify(body, null, 2))
      process.exit(1)
    }
  }
  console.log('OK: Phase 4 API smoke passed')
}

main().catch((e) => {
  console.error(e)
  process.exit(1)
})
