import { expect, test } from '@playwright/test'

const apiBase = process.env.API_BASE_URL || 'http://localhost:8080'

test.describe('Phase 0 smoke', () => {
  test('backend health returns ok', async ({ request }) => {
    const res = await request.get(`${apiBase}/health`)
    expect(res.ok()).toBeTruthy()

    const body = await res.json()
    expect(body.code).toBe(200)
    expect(body.data.status).toBe('ok')
    expect(body.data.db).toBe('connected')
  })

  test('frontend home page loads', async ({ page }) => {
    await page.goto('/')
    await expect(page.locator('h1')).toContainText('CRM')
    await expect(page.locator('h2')).toBeVisible()
  })
})
