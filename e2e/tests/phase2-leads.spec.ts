import { expect, test } from '@playwright/test'
import { apiBase, demoLeadId, loginDemoAdmin, loginDemoAdminUI } from './helpers/demo-session'

test.describe('Phase 2 leads E2E', () => {
  test('leads list loads and create lead appears in table', async ({ page }) => {
    const title = `E2E Lead ${Date.now()}`
    await loginDemoAdminUI(page)

    await page.goto('/leads')
    await expect(page.getByTestId('leads-page')).toBeVisible()
    await page.waitForResponse(
      (r) => r.url().includes('/api/leads') && r.request().method() === 'GET' && r.ok(),
      { timeout: 15_000 },
    )
    await expect(page.getByTestId('leads-list-table')).toBeVisible({ timeout: 15_000 })

    await page.getByTestId('lead-create-btn').click()
    const modalInput = page.locator('[role="dialog"][data-headlessui-state="open"] input').first()
    await expect(modalInput).toBeVisible({ timeout: 10_000 })
    await modalInput.fill(title)
    await page.locator('[role="dialog"][data-headlessui-state="open"]').getByRole('button', { name: /保存|Save/i }).click()

    await expect(page.getByTestId('leads-list-table')).toContainText(title, { timeout: 15_000 })
  })

  test('lead detail shows timeline and AI panel in preview mode', async ({ page }) => {
    await loginDemoAdminUI(page)
    await page.goto(`/leads/${demoLeadId}`)
    await expect(page.getByTestId('lead-detail-page')).toBeVisible({ timeout: 15_000 })
    await expect(page.getByTestId('activity-timeline')).toBeVisible()
    await expect(page.getByTestId('ai-relation-panel')).toBeVisible()
    await expect(page.getByTestId('ai-relation-panel')).not.toBeEmpty()
  })

  test('leads reports tab is present', async ({ page }) => {
    await loginDemoAdminUI(page)
    await page.goto('/leads')
    await expect(page.getByTestId('leads-main-tabs')).toBeVisible()
    await expect(page.getByRole('tab', { name: /^报表$|^Reports$/i })).toBeVisible()
  })

  test('leads stats API smoke', async ({ request }) => {
    const { headers } = await loginDemoAdmin(request)
    for (const path of ['/api/leads/stats/by-status', '/api/leads/stats/by-source', '/api/leads/stats/trend', '/api/leads/stats/funnel']) {
      const res = await request.get(`${apiBase}${path}`, { headers })
      expect(res.ok(), `${path} failed`).toBeTruthy()
    }
  })

  test('leads API CRUD smoke', async ({ request }) => {
    const ts = Date.now()
    const { headers } = await loginDemoAdmin(request)

    const createRes = await request.post(`${apiBase}/api/leads`, {
      headers,
      data: { title: `API Lead ${ts}`, source: 'website' },
    })
    expect(createRes.ok()).toBeTruthy()
    const leadId = (await createRes.json()).data.id as string

    const getRes = await request.get(`${apiBase}/api/leads/${leadId}`, { headers })
    expect(getRes.ok()).toBeTruthy()

    const patchRes = await request.patch(`${apiBase}/api/leads/${leadId}`, {
      headers,
      data: { status: 'contacted' },
    })
    expect(patchRes.ok()).toBeTruthy()

    const listRes = await request.get(`${apiBase}/api/leads?search=${encodeURIComponent(`API Lead ${ts}`)}`, {
      headers,
    })
    expect(listRes.ok()).toBeTruthy()
    const items = (await listRes.json()).data.items as Array<{ id: string }>
    expect(items.some((l) => l.id === leadId)).toBeTruthy()

    const delRes = await request.delete(`${apiBase}/api/leads/${leadId}`, { headers })
    expect(delRes.ok()).toBeTruthy()
  })
})

test.describe('Phase 2 leads preview path', () => {
  test('demo lead detail shows preview AI panel with copilot mock', async ({ page }) => {
    await loginDemoAdminUI(page)
    await page.goto(`/leads/${demoLeadId}?preview=1`)
    await expect(page.getByTestId('lead-detail-page')).toBeVisible({ timeout: 15_000 })
    await expect(page.getByTestId('ai-relation-panel')).toBeVisible()
    await expect(page.getByTestId('ai-relation-panel').getByTestId('ai-preview-badge')).toBeVisible()
    await expect(page.getByTestId('ai-preview-disclaimer')).toBeVisible()
    await expect(page.getByTestId('ai-churn-risk-bar')).toContainText('72%')
    await expect(page.getByTestId('ai-relation-panel')).toContainText('流失风险 72%')

    await page.getByTestId('ai-copilot-followup-btn').click()
    await expect(page.getByTestId('ai-copilot-output')).toBeVisible()
    await expect(page.getByTestId('ai-copilot-output')).not.toBeEmpty()

    await page.getByTestId('insight-adopt-btn').click()
    await expect(page.getByTestId('ai-relation-panel')).toContainText('已采纳演示话术', { timeout: 5_000 })
  })
})
