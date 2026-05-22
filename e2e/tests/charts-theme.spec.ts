import { test, expect } from '@playwright/test'

test.describe('charts theme toggle', () => {
  test('v1 to v2 updates document theme and chart colors', async ({ page }) => {
    const consoleErrors: string[] = []
    page.on('console', (msg) => {
      if (msg.type() === 'error') consoleErrors.push(msg.text())
    })
    page.on('pageerror', (err) => consoleErrors.push(err.message))

    await page.goto('/charts', { waitUntil: 'networkidle' })

    const htmlTheme = () => page.locator('html').getAttribute('data-theme')
    await expect.poll(htmlTheme).toBe('v1')

    const primaryBefore = await page.evaluate(() =>
      getComputedStyle(document.documentElement).getPropertyValue('--ds-chart-line-end').trim(),
    )

    await page.getByRole('button', { name: 'V2' }).click()
    await expect.poll(htmlTheme).toBe('v2')

    const primaryAfter = await page.evaluate(() =>
      getComputedStyle(document.documentElement).getPropertyValue('--ds-chart-line-end').trim(),
    )

    expect(primaryAfter).not.toBe(primaryBefore)

    const injectionWarns = consoleErrors.filter((t) =>
      /injection|UI_KIT_THEME|crm-ui-kit-theme/i.test(t),
    )
    expect(injectionWarns).toEqual([])

    await page.getByRole('button', { name: 'V1' }).click()
    await expect.poll(htmlTheme).toBe('v1')
  })
})
