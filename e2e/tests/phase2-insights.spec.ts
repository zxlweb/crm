import { expect, test } from '@playwright/test'
import { apiBase, loginDemoAdmin } from './helpers/demo-session'

test.describe('Phase 2 insights evaluate', () => {
  test('contact insights returns INS-001 for silent low engagement', async ({ request }) => {
    const ts = Date.now()
    const { headers } = await loginDemoAdmin(request)

    const createRes = await request.post(`${apiBase}/api/contacts`, {
      headers,
      data: {
        first_name: '洞察',
        last_name: `QA${ts}`,
        email: `insight-${ts}@example.com`,
      },
    })
    expect(createRes.ok()).toBeTruthy()
    const contactId = (await createRes.json()).data.id as string

    const evalRes = await request.post(`${apiBase}/api/contacts/${contactId}/insights/evaluate`, {
      headers,
      data: {},
    })
    expect(evalRes.ok()).toBeTruthy()
    const body = await evalRes.json()
    const items = body.data.items as Array<{ rule_id: string }>
    expect(items.some((i) => i.rule_id === 'INS-001')).toBeTruthy()
    expect(body.data.engagement_score).toBe(0)
  })

  test('account insights returns INS-001 for new account without activity', async ({ request }) => {
    const { headers } = await loginDemoAdmin(request)

    const accRes = await request.post(`${apiBase}/api/accounts`, {
      headers,
      data: { name: `Insight Account ${Date.now()}` },
    })
    expect(accRes.ok()).toBeTruthy()
    const accountId = (await accRes.json()).data.id as string

    const evalRes = await request.post(`${apiBase}/api/accounts/${accountId}/insights/evaluate`, {
      headers,
      data: {},
    })
    expect(evalRes.ok()).toBeTruthy()
    const items = (await evalRes.json()).data.items as Array<{ rule_id: string }>
    expect(items.some((i) => i.rule_id === 'INS-001')).toBeTruthy()
  })
})
