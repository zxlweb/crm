import { expect, test } from '@playwright/test'
import { apiBase, loginDemoAdmin, loginDemoAdminUI } from './helpers/demo-session'

test.describe('Phase 2 emotion journey API', () => {
  test('returns empty points for lead without activities', async ({ request }) => {
    const { headers } = await loginDemoAdmin(request)
    const ts = Date.now()

    const leadRes = await request.post(`${apiBase}/api/leads`, {
      headers,
      data: { title: `E2E Empty Journey ${ts}` },
    })
    expect(leadRes.ok()).toBeTruthy()
    const leadId = (await leadRes.json()).data.id as string

    const res = await request.get(`${apiBase}/api/leads/${leadId}/emotion-journey?range=90d`, { headers })
    expect(res.ok()).toBeTruthy()
    const data = (await res.json()).data
    expect(data.points).toEqual([])
    expect(data.summary.trend).toBe('flat')
  })

  test('returns points when activities have sentiment', async ({ request }) => {
    const { headers } = await loginDemoAdmin(request)
    const ts = Date.now()

    const leadRes = await request.post(`${apiBase}/api/leads`, {
      headers,
      data: { title: `E2E Journey Points ${ts}` },
    })
    expect(leadRes.ok()).toBeTruthy()
    const leadId = (await leadRes.json()).data.id as string

    const actRes = await request.post(`${apiBase}/api/activities`, {
      headers,
      data: {
        subject_type: 'lead',
        subject_id: leadId,
        event_type: 'call',
        direction: 'outbound',
        body: `Journey body ${ts}`,
        sentiment: 'positive',
        sentiment_source: 'manual',
      },
    })
    expect(actRes.ok()).toBeTruthy()

    const res = await request.get(`${apiBase}/api/leads/${leadId}/emotion-journey?range=all`, { headers })
    expect(res.ok()).toBeTruthy()
    const points = (await res.json()).data.points as Array<{ sentiment: string; sentiment_score: number }>
    expect(points.length).toBeGreaterThan(0)
    expect(points[0].sentiment).toBe('positive')
    expect(points[0].sentiment_score).toBe(2)
  })

  test('converted lead includes milestone', async ({ request }) => {
    const { headers } = await loginDemoAdmin(request)
    const ts = Date.now()

    const leadRes = await request.post(`${apiBase}/api/leads`, {
      headers,
      data: { title: `E2E Convert Journey ${ts}`, status: 'qualified' },
    })
    expect(leadRes.ok()).toBeTruthy()
    const leadId = (await leadRes.json()).data.id as string

    const convertRes = await request.post(`${apiBase}/api/leads/${leadId}/convert`, {
      headers,
      data: { create_account: { name: `Account ${ts}` } },
    })
    expect(convertRes.ok()).toBeTruthy()

    const res = await request.get(`${apiBase}/api/leads/${leadId}/emotion-journey`, { headers })
    expect(res.ok()).toBeTruthy()
    const milestones = (await res.json()).data.milestones as Array<{ type: string }>
    expect(milestones.some((m) => m.type === 'converted')).toBeTruthy()
  })
})

test.describe('Phase 2 emotion journey UI', () => {
  test('lead detail shows empty map then chart after activity', async ({ page, request }) => {
    const ts = Date.now()
    const bodyText = `E2E emotion map ${ts}`
    const { headers } = await loginDemoAdmin(request)

    const leadRes = await request.post(`${apiBase}/api/leads`, {
      headers,
      data: { title: `E2E UI Journey ${ts}` },
    })
    expect(leadRes.ok()).toBeTruthy()
    const leadId = (await leadRes.json()).data.id as string

    await loginDemoAdminUI(page)
    await page.goto(`/leads/${leadId}`)
    await expect(page.getByTestId('lead-detail-page')).toBeVisible()
    await expect(page.getByTestId('lead-decision-panel')).toBeVisible()

    const map = page.getByTestId('emotion-journey-map')
    await expect(map).toBeVisible({ timeout: 15_000 })
    await expect(map.getByText('暂无情绪触点')).toBeVisible()

    await page.getByTestId('activity-create-btn').click()
    await page.getByTestId('activity-form-body').fill(bodyText)
    const journeyReload = page.waitForResponse(
      (r) => r.url().includes('/emotion-journey') && r.request().method() === 'GET' && r.ok(),
    )
    await page.getByTestId('activity-form-submit').click()

    await expect(page.getByTestId('activity-timeline')).toContainText(bodyText, { timeout: 15_000 })
    await journeyReload
    await expect(map.getByText('暂无情绪触点')).toHaveCount(0, { timeout: 15_000 })
    await expect(map.getByText('趋势：')).toBeVisible()
  })

  test('account emotion tab shows journey map', async ({ page, request }) => {
    const { headers } = await loginDemoAdmin(request)
    const ts = Date.now()

    const accRes = await request.post(`${apiBase}/api/accounts`, {
      headers,
      data: { name: `E2E Emotion Account ${ts}` },
    })
    expect(accRes.ok()).toBeTruthy()
    const accountId = (await accRes.json()).data.id as string

    await loginDemoAdminUI(page)
    await page.goto(`/accounts/${accountId}`)
    await expect(page.getByTestId('account-detail-page')).toBeVisible()

    await page.getByRole('tab', { name: /情绪旅程|Emotion journey/i }).click()
    await expect(page.getByTestId('tab-emotion-journey')).toBeVisible()
    await expect(page.getByTestId('emotion-journey-map')).toBeVisible()
    await expect(page.getByTestId('emotion-journey-map').getByText('暂无情绪触点')).toBeVisible()
  })

  test('contact detail shows emotion journey section', async ({ page, request }) => {
    const { headers } = await loginDemoAdmin(request)
    const ts = Date.now()

    const contactRes = await request.post(`${apiBase}/api/contacts`, {
      headers,
      data: {
        first_name: '情绪',
        last_name: `E2E${ts}`,
        email: `emotion-${ts}@example.com`,
      },
    })
    expect(contactRes.ok()).toBeTruthy()
    const contactId = (await contactRes.json()).data.id as string

    await loginDemoAdminUI(page)
    await page.goto(`/contacts/${contactId}`)
    await expect(page.getByTestId('contact-detail-page')).toBeVisible()
    await expect(page.getByTestId('tab-emotion-journey')).toBeVisible()
    await expect(page.getByTestId('emotion-journey-map')).toBeVisible()
  })

  test('emotion range tabs are visible on lead detail', async ({ page, request }) => {
    const { headers } = await loginDemoAdmin(request)
    const ts = Date.now()

    const leadRes = await request.post(`${apiBase}/api/leads`, {
      headers,
      data: { title: `E2E Range ${ts}` },
    })
    expect(leadRes.ok()).toBeTruthy()
    const leadId = (await leadRes.json()).data.id as string

    await loginDemoAdminUI(page)
    await page.goto(`/leads/${leadId}`)
    await expect(page.getByTestId('emotion-journey-range')).toBeVisible()
    await expect(page.getByRole('tab', { name: /近30天|30 days/i })).toBeVisible()
    await expect(page.getByRole('tab', { name: /近90天|90 days/i })).toBeVisible()
    await expect(page.getByRole('tab', { name: /全部|All/i })).toBeVisible()
  })
})
