import { defineConfig, devices } from '@playwright/test';

// Playwright config for mockup screenshots
// Run with: npx playwright test --config=playwright.mockup.config.ts test/e2e/screenshots.spec.ts

export default defineConfig({
  testDir: './test/e2e',
  fullyParallel: false,
  forbidOnly: !!process.env.CI,
  retries: 0,
  workers: 1,
  reporter: 'list',
  use: {
    baseURL: 'http://localhost:8080',
    trace: 'off',
  },
  projects: [
    {
      name: 'chromium',
      use: { ...devices['Desktop Chrome'] },
    },
  ],
  webServer: {
    command: 'go run ./cmd/mockup',
    url: 'http://localhost:8080/health',
    reuseExistingServer: !process.env.CI,
    timeout: 30000,
  },
});
