import { expect, type APIRequestContext, type Page } from '@playwright/test'

export const apiBase = process.env.API_BASE_URL || 'http://localhost:8080'
export const demoTenantId = 'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa'
export const demoEmail = 'admin@demo.com'
export const demoPassword = 'password123'
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
