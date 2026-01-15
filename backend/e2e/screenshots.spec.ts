import { test, expect } from '@playwright/test';

// Screenshot test for UI mockup documentation
// Run with: npx playwright test --config=playwright.mockup.config.ts e2e/screenshots.spec.ts

test.describe('Admin UI Mockup Screenshots', () => {
  test.beforeEach(async ({ page }) => {
    // Set viewport for consistent screenshots
    await page.setViewportSize({ width: 1280, height: 800 });
  });

  test('Admin Login Page', async ({ page }) => {
    const response = await page.goto('/admin/login');
    expect(response?.status()).toBe(200);
    await page.screenshot({ path: '../docs/screenshots/admin-login.png', fullPage: true });
  });

  test('Admin Dashboard', async ({ page }) => {
    const response = await page.goto('/admin');
    expect(response?.status()).toBe(200);
    await page.screenshot({ path: '../docs/screenshots/admin-dashboard.png', fullPage: true });
  });

  test('Consultant Dashboard', async ({ page }) => {
    const response = await page.goto('/admin/my-dashboard');
    expect(response?.status()).toBe(200);
    await page.screenshot({ path: '../docs/screenshots/consultant-dashboard.png', fullPage: true });
  });

  test('Candidates List', async ({ page }) => {
    const response = await page.goto('/admin/candidates');
    expect(response?.status()).toBe(200);
    await page.screenshot({ path: '../docs/screenshots/admin-candidates.png', fullPage: true });
  });

  test('Candidates Filtered', async ({ page }) => {
    const response = await page.goto('/admin/candidates?status=prospecting');
    expect(response?.status()).toBe(200);
    await page.screenshot({ path: '../docs/screenshots/admin-candidates-filtered.png', fullPage: true });
  });

  test('Candidate Detail', async ({ page }) => {
    const response = await page.goto('/admin/candidates/2');
    expect(response?.status()).toBe(200);
    await page.screenshot({ path: '../docs/screenshots/admin-candidate-detail.png', fullPage: true });
  });

  test('Interaction Form', async ({ page }) => {
    const response = await page.goto('/admin/candidates/2/interaction');
    expect(response?.status()).toBe(200);
    await page.screenshot({ path: '../docs/screenshots/interaction-form.png', fullPage: true });
  });

  test('Consultant Performance Report', async ({ page }) => {
    const response = await page.goto('/admin/reports/consultants');
    expect(response?.status()).toBe(200);
    await page.screenshot({ path: '../docs/screenshots/consultant-report.png', fullPage: true });
  });

  test('Document Review', async ({ page }) => {
    const response = await page.goto('/admin/documents');
    expect(response?.status()).toBe(200);
    await page.screenshot({ path: '../docs/screenshots/document-review.png', fullPage: true });
  });

  test('Portal Documents', async ({ page }) => {
    const response = await page.goto('/portal/documents');
    expect(response?.status()).toBe(200);
    await page.screenshot({ path: '../docs/screenshots/portal-documents.png', fullPage: true });
  });
});
