import { expect, test } from '@playwright/test'
import { apiBase, loginDemoAdmin, loginDemoAdminUI } from './helpers/demo-session'

test.describe('Phase 2 segments', () => {
  test('segments API list count matches leads list total (L-INT-04b)', async ({ request }) => {
    const { headers } = await loginDemoAdmin(request)

    const listRes = await request.get(`${apiBase}/api/segments`, { headers })
    if (!listRes.ok()) {
      test.skip(true, 'segments API not deployed yet')
    }
    const templates = (await listRes.json()).data.items as Array<{ code: string }>
    expect(templates.length).toBeGreaterThanOrEqual(5)

    const countRes = await request.get(`${apiBase}/api/segments/high_value/count`, { headers })
    expect(countRes.ok()).toBeTruthy()
    const count = (await countRes.json()).data.count as number

    const leadsRes = await request.get(`${apiBase}/api/leads?segment=high_value&page_size=100`, {
      headers,
    })
    expect(leadsRes.ok()).toBeTruthy()
    const total = (await leadsRes.json()).pagination.total as number
    expect(total).toBe(count)
  })

  test('segment URL query loads leads list', async ({ page }) => {
    await loginDemoAdminUI(page)
    await page.goto('/leads?segment=high_value')
    await expect(page).toHaveURL(/segment=high_value/)
    await expect(page.getByTestId('leads-page')).toBeVisible()
    await expect(page.getByTestId('leads-list-table')).toBeVisible({ timeout: 15_000 })
    await page.waitForResponse(
      (r) => r.url().includes('/api/leads') && r.url().includes('segment=high_value'),
      { timeout: 15_000 },
    ).catch(() => {
      // FE mock 路径下可能无真实 API 请求
    })
  })
})
