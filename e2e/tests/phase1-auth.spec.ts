import { expect, test } from '@playwright/test'

const apiBase = process.env.API_BASE_URL || 'http://localhost:8080'

test.describe('Phase 1 auth', () => {
  test('register new tenant admin lands on home', async ({ page, request }) => {
    const ts = Date.now()
    const email = `e2e-reg-${ts}@example.com`

    await page.goto('/login?mode=register')
    await expect(page.getByTestId('auth-form-title')).toContainText(/工作区|workspace/i)

    await page.getByTestId('auth-input-name').fill('E2E User')
    await page.getByTestId('auth-input-company').fill(`E2E Co ${ts}`)
    await page.getByTestId('auth-input-email').fill(email)
    await page.getByTestId('auth-input-password').fill('password123')
    await page.getByTestId('auth-submit').click()

    await expect(page).toHaveURL(/\//)
    await expect(page).not.toHaveURL(/\/login/)

    const res = await request.post(`${apiBase}/api/auth/login`, {
      data: { email, password: 'password123' },
    })
    expect(res.ok()).toBeTruthy()
    const body = await res.json()
    expect(body.code).toBe(200)
    expect(body.data.user.email).toBe(email)
  })

  test('demo super admin login lands on admin', async ({ page }) => {
    await page.goto('/login')
    await page.getByTestId('auth-input-email').fill('admin@demo.com')
    await page.getByTestId('auth-input-password').fill('password123')
    await page.getByTestId('auth-submit').click()
    await expect(page).toHaveURL(/\/admin/)
  })

  test('admin tenants nav scrolls without hash in URL', async ({ page }) => {
    await page.goto('/login')
    await page.getByTestId('auth-input-email').fill('admin@demo.com')
    await page.getByTestId('auth-input-password').fill('password123')
    await page.getByTestId('auth-submit').click()
    await expect(page).toHaveURL(/\/admin\/?$/)

    await page.getByTestId('admin-nav-tenants').click()
    await expect(page).toHaveURL(/\/admin\/?$/)
    await expect(page).not.toHaveURL(/#/)
    await expect(page.locator('#tenants')).toBeInViewport()
  })
})
