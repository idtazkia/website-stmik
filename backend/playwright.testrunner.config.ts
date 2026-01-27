import { defineConfig, devices } from '@playwright/test';

// Playwright config for testrunner - uses externally managed server
// The testrunner starts the server and passes BASE_URL via environment variable

export default defineConfig({
  testDir: './test/e2e',
  fullyParallel: true,
  forbidOnly: true,
  retries: 2,
  workers: 1, // Sequential for stability with shared database
  reporter: [['list'], ['html', { open: 'never' }]],
  use: {
    baseURL: process.env.BASE_URL || 'http://localhost:8080',
    trace: 'on-first-retry',
  },
  projects: [
    {
      name: 'chromium',
      use: { ...devices['Desktop Chrome'] },
    },
  ],
  // No webServer - testrunner manages the server externally
});
