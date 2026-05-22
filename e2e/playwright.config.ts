import { defineConfig, devices } from '@playwright/test'

const apiBase = process.env.API_BASE_URL || 'http://localhost:8080'
const webBase = process.env.WEB_BASE_URL || 'http://localhost:3000'

export default defineConfig({
  testDir: './tests',
  timeout: 30_000,
  retries: 0,
  use: {
    baseURL: webBase,
    trace: 'on-first-retry',
  },
  projects: [
    { name: 'chromium', use: { ...devices['Desktop Chrome'] } },
  ],
  webServer: [
    {
      command: 'cd ../backend && go run ./cmd/api/',
      url: `${apiBase}/health`,
      reuseExistingServer: true,
      timeout: 120_000,
    },
    {
      command: 'cd ../frontend && npm run dev',
      url: webBase,
      reuseExistingServer: true,
      timeout: 120_000,
    },
  ],
})
