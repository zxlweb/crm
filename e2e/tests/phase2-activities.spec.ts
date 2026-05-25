import { expect, test } from '@playwright/test'
import { apiBase, loginDemoAdmin, loginDemoAdminUI } from './helpers/demo-session'

test.describe('Phase 2 activities', () => {
  test('create activity shows on lead timeline', async ({ page, request }) => {
    const ts = Date.now()
    const bodyText = `E2E call notes ${ts}`
    const { headers } = await loginDemoAdmin(request)

    const leadRes = await request.post(`${apiBase}/api/leads`, {
      headers,
      data: { title: `E2E Lead ${ts}`, source: 'website' },
    })
    expect(leadRes.ok()).toBeTruthy()
    const leadId = (await leadRes.json()).data.id as string

    await loginDemoAdminUI(page)

    await page.goto(`/leads/${leadId}`)
    await expect(page.getByTestId('lead-detail-page')).toBeVisible()
    await expect(page.getByTestId('activity-timeline')).toBeVisible()

    await page.getByTestId('activity-create-btn').click()
    await expect(page.getByTestId('activity-form')).toBeVisible()
    await page.getByTestId('activity-form-body').fill(bodyText)
    await page.getByTestId('activity-form-submit').click()

    await expect(page.getByTestId('activity-timeline')).toContainText(bodyText, { timeout: 15_000 })
  })

  test('activity form exposes sentiment field', async ({ page, request }) => {
    const ts = Date.now()
    const { headers } = await loginDemoAdmin(request)

    const leadRes = await request.post(`${apiBase}/api/leads`, {
      headers,
      data: { title: `E2E Sentiment Form ${ts}` },
    })
    expect(leadRes.ok()).toBeTruthy()
    const leadId = (await leadRes.json()).data.id as string

    await loginDemoAdminUI(page)
    await page.goto(`/leads/${leadId}`)
    await page.getByTestId('activity-create-btn').click()
    await expect(page.getByTestId('activity-form-sentiment')).toBeVisible()
  })

  test('timeline shows sentiment label for API-created activity', async ({ page, request }) => {
    const ts = Date.now()
    const bodyText = `E2E sentiment timeline ${ts}`
    const { headers } = await loginDemoAdmin(request)

    const leadRes = await request.post(`${apiBase}/api/leads`, {
      headers,
      data: { title: `E2E Sentiment Timeline ${ts}` },
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
        body: bodyText,
        sentiment: 'positive',
        sentiment_source: 'manual',
      },
    })
    expect(actRes.ok()).toBeTruthy()

    await loginDemoAdminUI(page)
    await page.goto(`/leads/${leadId}`)
    await expect(page.getByTestId('activity-timeline')).toContainText(bodyText, { timeout: 15_000 })
    await expect(page.getByTestId('activity-timeline')).toContainText('积极', { timeout: 15_000 })
  })

  test('activities API CRUD via request', async ({ request }) => {
    const ts = Date.now()
    const { headers } = await loginDemoAdmin(request)

    const leadRes = await request.post(`${apiBase}/api/leads`, { headers, data: { title: `Lead ${ts}` } })
    expect(leadRes.ok()).toBeTruthy()
    const leadId = (await leadRes.json()).data.id as string

    const createRes = await request.post(`${apiBase}/api/activities`, {
      headers,
      data: {
        subject_type: 'lead',
        subject_id: leadId,
        event_type: 'call',
        direction: 'outbound',
        body: `API body ${ts}`,
        sentiment: 'positive',
        sentiment_source: 'manual',
      },
    })
    expect(createRes.ok()).toBeTruthy()
    const activityId = (await createRes.json()).data.id as string

    const listRes = await request.get(
      `${apiBase}/api/activities?subject_type=lead&subject_id=${leadId}`,
      { headers },
    )
    expect(listRes.ok()).toBeTruthy()
    const listBody = await listRes.json()
    expect(listBody.data.items.some((a: { id: string }) => a.id === activityId)).toBeTruthy()

    const summaryRes = await request.get(`${apiBase}/api/activities/summary?subject_type=lead&subject_id=${leadId}`, {
      headers,
    })
    expect(summaryRes.ok()).toBeTruthy()
    expect((await summaryRes.json()).data.total).toBeGreaterThan(0)

    const delRes = await request.delete(`${apiBase}/api/activities/${activityId}`, { headers })
    expect(delRes.ok()).toBeTruthy()
  })
})
