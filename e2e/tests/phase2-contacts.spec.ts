import { expect, test } from '@playwright/test'

const apiBase = process.env.API_BASE_URL || 'http://localhost:8080'
const demoTenantId = 'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa'
const demoEmail = 'admin@demo.com'
const demoPassword = 'password123'

async function loginDemoAdmin(request: import('@playwright/test').APIRequestContext) {
  const res = await request.post(`${apiBase}/api/auth/login`, {
    data: { email: demoEmail, password: demoPassword },
  })
  expect(res.ok()).toBeTruthy()
  const body = await res.json()
  return {
    token: body.data.access_token as string,
    headers: {
      Authorization: `Bearer ${body.data.access_token}`,
      'X-Tenant-ID': demoTenantId,
      'Content-Type': 'application/json',
    },
  }
}

test.describe('Phase 2 contacts', () => {
  test('create contact shows in list table', async ({ page, request }) => {
    const ts = Date.now()
    const email = `e2e-ui-${ts}@example.com`
    const { headers } = await loginDemoAdmin(request)

    await page.goto('/login')
    await page.getByTestId('auth-input-email').fill(demoEmail)
    await page.getByTestId('auth-input-password').fill(demoPassword)
    await page.getByTestId('auth-submit').click()
    await expect(page).toHaveURL(/\/admin/)

    await page.context().addCookies([
      {
        name: 'crm.tenant_id',
        value: demoTenantId,
        domain: 'localhost',
        path: '/',
      },
    ])

    await page.goto('/contacts')
    await expect(page.getByTestId('contacts-page')).toBeVisible()
    await page.waitForResponse(
      (r) => r.url().includes('/api/contacts') && r.request().method() === 'GET' && r.ok(),
      { timeout: 15_000 },
    )
    await expect(page.getByTestId('contacts-list-table')).toBeVisible({ timeout: 15_000 })

    await page.getByTestId('contact-create-btn').click()
    await expect(page.getByTestId('contact-form')).toBeVisible()
    await page.getByTestId('contact-form-first-name').fill('E2E')
    await page.getByTestId('contact-form-last-name').fill(`UI${ts}`)
    await page.getByTestId('contact-form-email').fill(email)
    await page.getByTestId('contact-form-submit').click()

    await expect(page.getByTestId('contacts-list-table')).toContainText(`E2E UI${ts}`, {
      timeout: 15_000,
    })

    const listRes = await request.get(`${apiBase}/api/contacts?search=${encodeURIComponent(email)}`, {
      headers,
    })
    expect(listRes.ok()).toBeTruthy()
    const items = (await listRes.json()).data.items as Array<{ email: string }>
    expect(items.some((c) => c.email === email)).toBeTruthy()
  })

  test('contacts API CRUD and account link', async ({ request }) => {
    const { headers } = await loginDemoAdmin(request)

    const accRes = await request.post(`${apiBase}/api/accounts`, {
      headers,
      data: { name: `E2E Account ${Date.now()}` },
    })
    expect(accRes.ok()).toBeTruthy()
    const accountId = (await accRes.json()).data.id as string

    const createRes = await request.post(`${apiBase}/api/contacts`, {
      headers,
      data: {
        account_id: accountId,
        first_name: 'E2E',
        last_name: 'Contact',
        email: `e2e-${Date.now()}@example.com`,
        is_primary: true,
      },
    })
    expect(createRes.ok()).toBeTruthy()
    const contactId = (await createRes.json()).data.id as string

    const listRes = await request.get(`${apiBase}/api/accounts/${accountId}/contacts`, { headers })
    expect(listRes.ok()).toBeTruthy()
    const items = (await listRes.json()).data.items as Array<{ id: string }>
    expect(items.some((c) => c.id === contactId)).toBeTruthy()

    const delRes = await request.delete(`${apiBase}/api/contacts/${contactId}`, { headers })
    expect(delRes.ok()).toBeTruthy()
  })
})
