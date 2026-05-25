import { expect, test } from '@playwright/test'
import { demoLeadId, loginDemoAdminUI } from './helpers/demo-session'

/**
 * PRD §15.2 老板汇报 — 3 分钟演示路径（全链路）
 * @see docs/prd/phase-2-relationship-crm-prd.md §15.2
 */
test.describe('Phase 2 §15.2 demo path', () => {
  test.describe.configure({ mode: 'serial' })

  test('full path without 5xx: login → 华创科技 → emotion → AI preview → reports', async ({
    page,
  }) => {
    const serverErrors: Array<{ url: string; status: number }> = []
    page.on('response', (res) => {
      if (res.status() >= 500) {
        serverErrors.push({ url: res.url(), status: res.status() })
      }
    })

    // Step 1 — 登录（种子账号 admin@demo.com，演示租户 cookie 在 helper 中附加）
    await loginDemoAdminUI(page)
    await expect(page).toHaveURL(/\/admin/)

    // Step 2 — 线索列表打开「⭐ 华创科技」，生命周期徽章 + 互动分
    await page.goto('/leads')
    await expect(page.getByTestId('leads-page')).toBeVisible()
    await page.waitForResponse(
      (r) => r.url().includes('/api/leads') && r.request().method() === 'GET' && r.ok(),
      { timeout: 15_000 },
    )
    await expect(page.getByTestId('leads-list-table')).toBeVisible({ timeout: 15_000 })

    const searchBox = page.getByRole('searchbox', { name: /搜索标题|Search title/i })
    const searchResp = page.waitForResponse(
      (r) =>
        r.url().includes('/api/leads') &&
        r.request().method() === 'GET' &&
        r.url().includes('search') &&
        r.ok(),
    )
    await searchBox.fill('华创')
    await searchBox.press('Enter')
    await searchResp

    const demoLink = page.getByRole('link', { name: /华创科技/ })
    await expect(demoLink.first()).toBeVisible({ timeout: 15_000 })
    await demoLink.first().click()

    await expect(page).toHaveURL(new RegExp(`/leads/${demoLeadId}`))
    await expect(page.getByTestId('lead-detail-page')).toBeVisible({ timeout: 15_000 })
    await expect(page.getByTestId('lead-detail-header')).toBeVisible()
    await expect(page.getByTestId('lead-lifecycle-badges')).toBeVisible()
    await expect(page.getByTestId('lead-detail-metrics')).toBeVisible()
    await expect(page.getByTestId('lead-detail-metrics')).toContainText(/\d+/)

    // Step 3 — 情绪旅程（详情内嵌于决策面板，非独立 Tab）
    await expect(page.getByTestId('lead-decision-panel')).toBeVisible()
    await expect(page.getByTestId('emotion-journey-range')).toBeVisible()
    const map = page.getByTestId('emotion-journey-map')
    await expect(map).toBeVisible({ timeout: 15_000 })
    await expect(map.getByText('暂无情绪触点')).toHaveCount(0, { timeout: 15_000 })
    await expect(map.getByText(/趋势：|Trend:/i)).toBeVisible()

    // Step 4–5 — AI 关系助手：流失风险样例、Preview 角标、Copilot Mock、采纳
    const aiPanel = page.getByTestId('ai-relation-panel')
    await expect(aiPanel).toBeVisible()
    await expect(aiPanel.getByTestId('ai-preview-badge')).toBeVisible()
    await expect(aiPanel.getByTestId('ai-preview-badge')).toContainText(/Preview|演示/i)
    await expect(aiPanel.getByTestId('ai-preview-disclaimer')).toContainText(/演示样例|Preview/i)
    await expect(page.getByTestId('ai-churn-risk-bar')).toBeVisible()
    await expect(page.getByTestId('ai-churn-risk-bar')).toContainText('72')
    await expect(aiPanel).toContainText(/流失风险 72%|Churn risk 72%/i)
    await expect(aiPanel).toContainText(/高意向冷却|Hot lead cooling/i)

    await page.getByTestId('ai-copilot-followup-btn').click()
    const copilotOut = page.getByTestId('ai-copilot-output')
    await expect(copilotOut).toBeVisible({ timeout: 5_000 })
    await expect(copilotOut).toContainText(/部署周期|milestones/i)

    await page.getByTestId('insight-adopt-btn').click()
    await expect(copilotOut).toBeVisible()

    // Step 6（可选）— Leads 报表 Tab：折线 + 环形 + 漏斗
    await page.goto('/leads')
    await expect(page.getByTestId('leads-main-tabs')).toBeVisible()
    await page.getByRole('tab', { name: /^报表$|^Reports$/i }).click()
    const reports = page.getByTestId('leads-tab-reports')
    await expect(reports).toBeVisible({ timeout: 15_000 })
    await page.waitForResponse(
      (r) => r.url().includes('/api/leads/stats') && r.ok(),
      { timeout: 15_000 },
    ).catch(() => undefined)
    await expect(reports.getByText(/按状态|By status|来源|Source/i).first()).toBeVisible()
    await expect(reports.getByText(/趋势|Trend/i).first()).toBeVisible()
    await expect(reports.getByText(/漏斗|Funnel/i).first()).toBeVisible()

    expect(serverErrors, `unexpected 5xx: ${JSON.stringify(serverErrors)}`).toEqual([])
  })

  test('demo lead direct URL with ?preview=1 keeps preview badge', async ({ page }) => {
    await loginDemoAdminUI(page)
    await page.goto(`/leads/${demoLeadId}?preview=1`)
    await expect(page.getByTestId('lead-detail-page')).toBeVisible({ timeout: 15_000 })
    await expect(page.getByTestId('ai-relation-panel').getByTestId('ai-preview-badge')).toBeVisible()
    await expect(page.getByTestId('emotion-journey-map')).toBeVisible({ timeout: 15_000 })
  })
})
