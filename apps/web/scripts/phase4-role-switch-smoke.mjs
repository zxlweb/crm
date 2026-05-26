/**
 * Phase 4 — 多角色切换 API smoke
 * 前置：后端已启动且迁移含 00017（multi-role@demo.com）
 * 运行：node apps/web/scripts/phase4-role-switch-smoke.mjs
 */
const base = process.env.API_BASE || 'http://localhost:8080'
const tenant = 'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa'
const email = 'multi-role@demo.com'
const password = 'password123'

function tenantHeaders(token) {
  return {
    Authorization: `Bearer ${token}`,
    'X-Tenant-ID': tenant,
    'Content-Type': 'application/json',
  }
}

async function post(path, body, headers = {}) {
  const res = await fetch(`${base}${path}`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json', ...headers },
    body: JSON.stringify(body),
  })
  const data = await res.json()
  return { res, data }
}

async function get(path, headers) {
  const res = await fetch(`${base}${path}`, { headers })
  const data = await res.json()
  return { res, data }
}

function assertOk(label, res, data) {
  if (!res.ok || data.code !== 200) {
    console.error(label, res.status, JSON.stringify(data, null, 2))
    process.exit(1)
  }
}

async function main() {
  const login = await post('/api/auth/login', { email, password })
  assertOk('login', login.res, login.data)
  let token = login.data.data.access_token

  const switched = await post(
    '/api/auth/switch-tenant',
    { tenant_id: tenant },
    { Authorization: `Bearer ${token}` },
  )
  assertOk('switch-tenant', switched.res, switched.data)
  token = switched.data.data.access_token
  const roles = switched.data.data.roles ?? []
  if (roles.length < 2) {
    console.error('expected >= 2 roles after switch-tenant, got', roles.length)
    process.exit(1)
  }
  const currentId = switched.data.data.current_role?.id
  const other = roles.find((r) => r.id !== currentId)
  if (!other) {
    console.error('no alternate role to switch to')
    process.exit(1)
  }
  console.log('switch-tenant roles:', roles.map((r) => r.name).join(', '), 'current:', currentId)

  const myRoles = await get('/api/rbac/my-roles', tenantHeaders(token))
  assertOk('my-roles', myRoles.res, myRoles.data)
  if ((myRoles.data.data?.length ?? 0) < 2) {
    console.error('my-roles length < 2')
    process.exit(1)
  }

  const roleSwitch = await post(
    '/api/auth/switch-role',
    { role_id: other.id },
    tenantHeaders(token),
  )
  assertOk('switch-role', roleSwitch.res, roleSwitch.data)
  token = roleSwitch.data.data.access_token
  if (roleSwitch.data.data.current_role?.id !== other.id) {
    console.error('current_role mismatch', roleSwitch.data.data.current_role, 'expected', other.id)
    process.exit(1)
  }
  console.log('switch-role ->', roleSwitch.data.data.current_role.name)

  const forbidden = await post(
    '/api/auth/switch-role',
    { role_id: '00000000-0000-4000-8000-000000000099' },
    tenantHeaders(token),
  )
  if (forbidden.res.status !== 403 && forbidden.data.code !== 403) {
    console.error('expected 403 for invalid role', forbidden.res.status, forbidden.data)
    process.exit(1)
  }

  const check = await post(
    '/api/rbac/check',
    { resource: 'settings', action: 'update' },
    tenantHeaders(token),
  )
  assertOk('rbac/check', check.res, check.data)
  if (check.data.data?.allowed !== false) {
    console.error('expected settings:update denied for Sales Manager / Viewer')
    process.exit(1)
  }

  console.log('OK: Phase 4 role-switch smoke passed')
}

main().catch((e) => {
  console.error(e)
  process.exit(1)
})
