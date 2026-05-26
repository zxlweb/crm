import { expect, type APIRequestContext, type Page } from '@playwright/test'

export const apiBase = process.env.API_BASE_URL || 'http://localhost:8080'
export const demoTenantId = 'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa'
export const demoEmail = 'admin@demo.com'
export const demoPassword = 'password123'
export const multiRoleEmail = 'multi-role@demo.com'
export const roleSalesManagerId = 'e1e1e1e1-e1e1-e1e1-e1e1-e1e1e1e1e1e1'
export const roleViewerId = 'e2e2e2e2-e2e2-e2e2-e2e2-e2e2e2e2e2e2'
export const demoLeadId = 'a1000000-0000-4000-8000-000000000001'

export async function loginDemoAdmin(request: APIRequestContext) {
  const res = await request.post(`${apiBase}/api/auth/login`, {
    data: { email: demoEmail, password: demoPassword },
  })
  expect(res.ok()).toBeTruthy()
  const body = await res.json()
  const token = body.data.access_token as string
  return {
    token,
    headers: {
      Authorization: `Bearer ${token}`,
      'X-Tenant-ID': demoTenantId,
      'Content-Type': 'application/json',
    },
  }
}

export async function loginDemoAdminUI(page: Page) {
  await page.goto('/login')
  await page.getByTestId('auth-input-email').fill(demoEmail)
  await page.getByTestId('auth-input-password').fill(demoPassword)
  await page.getByTestId('auth-submit').click()
  await expect(page).toHaveURL(/\/admin/)
  await attachDemoTenant(page)
}

export async function attachDemoTenant(page: Page) {
  await page.context().addCookies([
    {
      name: 'crm.tenant_id',
      value: demoTenantId,
      domain: 'localhost',
      path: '/',
    },
  ])
}

/** 多角色演示账号：登录 + 切到 Demo Corp（含 roles / JWT active_role） */
export async function loginMultiRole(request: APIRequestContext) {
  const loginRes = await request.post(`${apiBase}/api/auth/login`, {
    data: { email: multiRoleEmail, password: demoPassword },
  })
  expect(loginRes.ok()).toBeTruthy()
  const loginBody = await loginRes.json()
  let token = loginBody.data.access_token as string

  const switchRes = await request.post(`${apiBase}/api/auth/switch-tenant`, {
    data: { tenant_id: demoTenantId },
    headers: { Authorization: `Bearer ${token}` },
  })
  expect(switchRes.ok()).toBeTruthy()
  const switchBody = await switchRes.json()
  token = switchBody.data.access_token as string
  const roles = switchBody.data.roles as Array<{ id: string; name: string }>
  expect(roles.length).toBeGreaterThanOrEqual(2)

  return {
    token,
    roles,
    currentRoleId: switchBody.data.current_role?.id as string | undefined,
    headers: {
      Authorization: `Bearer ${token}`,
      'X-Tenant-ID': demoTenantId,
      'Content-Type': 'application/json',
    },
  }
}

export async function loginMultiRoleUI(page: Page) {
  await page.goto('/login')
  await page.getByTestId('auth-input-email').fill(multiRoleEmail)
  await page.getByTestId('auth-input-password').fill(demoPassword)
  await Promise.all([
    page.waitForURL((url) => !url.pathname.includes('/login'), { timeout: 25_000 }),
    page.getByTestId('auth-submit').click(),
  ])
  await page
    .waitForResponse((r) => r.url().includes('/api/auth/switch-tenant') && r.ok(), { timeout: 15_000 })
    .catch(() => undefined)
}
