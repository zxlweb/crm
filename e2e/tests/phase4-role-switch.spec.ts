import { expect, test } from '@playwright/test'
import {
  apiBase,
  demoTenantId,
  loginMultiRole,
  loginMultiRoleUI,
  multiRoleEmail,
  roleSalesManagerId,
  roleViewerId,
} from './helpers/demo-session'

test.describe('Phase 4 role switch', () => {
  test('API: switch-tenant exposes roles, switch-role updates current_role, invalid role 403', async ({
    request,
  }) => {
    const session = await loginMultiRole(request)
    const otherId =
      session.currentRoleId === roleViewerId ? roleSalesManagerId : roleViewerId

    const switchRes = await request.post(`${apiBase}/api/auth/switch-role`, {
      data: { role_id: otherId },
      headers: session.headers,
    })
    expect(switchRes.ok()).toBeTruthy()
    const switchBody = await switchRes.json()
    expect(switchBody.code).toBe(200)
    expect(switchBody.data.current_role.id).toBe(otherId)

    const badRes = await request.post(`${apiBase}/api/auth/switch-role`, {
      data: { role_id: '00000000-0000-4000-8000-000000000099' },
      headers: {
        ...session.headers,
        Authorization: `Bearer ${switchBody.data.access_token}`,
      },
    })
    expect(badRes.status()).toBe(403)

    const checkRes = await request.post(`${apiBase}/api/rbac/check`, {
      data: { resource: 'settings', action: 'update' },
      headers: {
        ...session.headers,
        Authorization: `Bearer ${switchBody.data.access_token}`,
      },
    })
    expect(checkRes.ok()).toBeTruthy()
    const checkBody = await checkRes.json()
    expect(checkBody.data.allowed).toBe(false)
  })

  test('UI: multi-role user can switch role via the topbar user menu', async ({ page }) => {
    test.setTimeout(60_000)
    await loginMultiRoleUI(page)

    const trigger = page.getByTestId('user-menu-trigger').first()
    await expect(trigger).toBeVisible({ timeout: 25_000 })
    await trigger.click()

    const rolesSection = page.getByTestId('user-menu-roles')
    await expect(rolesSection).toBeVisible({ timeout: 10_000 })

    const roleItems = rolesSection.getByTestId('user-menu-role-item')
    await expect(roleItems).toHaveCount(2)

    // Determine which role is currently active (aria-checked="true") and pick the other.
    const activeIndex = await roleItems.evaluateAll((nodes) =>
      nodes.findIndex((n) => n.getAttribute('aria-checked') === 'true'),
    )
    const currentRoleName = await roleItems.nth(Math.max(activeIndex, 0)).innerText()
    const targetId = currentRoleName.includes('Viewer') ? roleSalesManagerId : roleViewerId
    const targetIndex = activeIndex === 0 ? 1 : 0

    const switchResp = page.waitForResponse(
      (r) => r.url().includes('/api/auth/switch-role') && r.request().method() === 'POST',
    )
    await roleItems.nth(targetIndex).click()
    const res = await switchResp
    expect(res.ok()).toBeTruthy()
    const body = await res.json()
    expect(body.data.current_role.id).toBe(targetId)
  })

  test('login response lists multi-role account email', async ({ request }) => {
    const res = await request.post(`${apiBase}/api/auth/login`, {
      data: { email: multiRoleEmail, password: 'password123' },
    })
    expect(res.ok()).toBeTruthy()
    const body = await res.json()
    expect(body.data.user.email).toBe(multiRoleEmail)

    const switchRes = await request.post(`${apiBase}/api/auth/switch-tenant`, {
      data: { tenant_id: demoTenantId },
      headers: { Authorization: `Bearer ${body.data.access_token}` },
    })
    expect(switchRes.ok()).toBeTruthy()
    const switched = await switchRes.json()
    expect(switched.data.roles?.length ?? 0).toBeGreaterThanOrEqual(2)
  })
})
