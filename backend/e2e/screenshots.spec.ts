import { test, expect } from '@playwright/test';

// Screenshot test for UI mockup documentation
// Run with: npx playwright test --config=playwright.mockup.config.ts e2e/screenshots.spec.ts

test.describe('Admin UI Mockup Screenshots', () => {
  test.beforeEach(async ({ page }) => {
    // Set viewport for consistent screenshots
    await page.setViewportSize({ width: 1280, height: 800 });
  });

  // Auth
  test('Admin Login Page', async ({ page }) => {
    const response = await page.goto('/admin/login');
    expect(response?.status()).toBe(200);
    await page.screenshot({ path: '../docs/screenshots/admin-login.png', fullPage: true });
  });

  // Dashboards
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

  // Candidates
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

  // Documents
  test('Document Review', async ({ page }) => {
    const response = await page.goto('/admin/documents');
    expect(response?.status()).toBe(200);
    await page.screenshot({ path: '../docs/screenshots/document-review.png', fullPage: true });
  });

  // Marketing
  test('Campaigns', async ({ page }) => {
    const response = await page.goto('/admin/campaigns');
    expect(response?.status()).toBe(200);
    await page.screenshot({ path: '../docs/screenshots/settings-campaigns.png', fullPage: true });
  });

  test('Referrers', async ({ page }) => {
    const response = await page.goto('/admin/referrers');
    expect(response?.status()).toBe(200);
    await page.screenshot({ path: '../docs/screenshots/settings-referrers.png', fullPage: true });
  });

  test('Commissions', async ({ page }) => {
    const response = await page.goto('/admin/commissions');
    expect(response?.status()).toBe(200);
    await page.screenshot({ path: '../docs/screenshots/admin-commissions.png', fullPage: true });
  });

  // Reports
  test('Consultant Performance Report', async ({ page }) => {
    const response = await page.goto('/admin/reports/consultants');
    expect(response?.status()).toBe(200);
    await page.screenshot({ path: '../docs/screenshots/consultant-report.png', fullPage: true });
  });

  test('Funnel Report', async ({ page }) => {
    const response = await page.goto('/admin/reports/funnel');
    expect(response?.status()).toBe(200);
    await page.screenshot({ path: '../docs/screenshots/report-funnel.png', fullPage: true });
  });

  test('Campaign ROI Report', async ({ page }) => {
    const response = await page.goto('/admin/reports/campaigns');
    expect(response?.status()).toBe(200);
    await page.screenshot({ path: '../docs/screenshots/report-campaigns.png', fullPage: true });
  });

  // Settings
  test('Settings Users', async ({ page }) => {
    const response = await page.goto('/admin/settings/users');
    expect(response?.status()).toBe(200);
    await page.screenshot({ path: '../docs/screenshots/settings-users.png', fullPage: true });
  });

  test('Settings Programs', async ({ page }) => {
    const response = await page.goto('/admin/settings/programs');
    expect(response?.status()).toBe(200);
    await page.screenshot({ path: '../docs/screenshots/settings-programs.png', fullPage: true });
  });

  test('Settings Categories', async ({ page }) => {
    const response = await page.goto('/admin/settings/categories');
    expect(response?.status()).toBe(200);
    await page.screenshot({ path: '../docs/screenshots/settings-categories.png', fullPage: true });
  });
});

test.describe('Portal UI Mockup Screenshots', () => {
  test.beforeEach(async ({ page }) => {
    await page.setViewportSize({ width: 1280, height: 800 });
  });

  test('Portal Dashboard', async ({ page }) => {
    const response = await page.goto('/portal');
    expect(response?.status()).toBe(200);
    await page.screenshot({ path: '../docs/screenshots/portal-dashboard.png', fullPage: true });
  });

  test('Portal Registration', async ({ page }) => {
    const response = await page.goto('/portal/register');
    expect(response?.status()).toBe(200);
    await page.screenshot({ path: '../docs/screenshots/portal-registration.png', fullPage: true });
  });

  test('Portal Documents', async ({ page }) => {
    const response = await page.goto('/portal/documents');
    expect(response?.status()).toBe(200);
    await page.screenshot({ path: '../docs/screenshots/portal-documents.png', fullPage: true });
  });

  test('Portal Payments', async ({ page }) => {
    const response = await page.goto('/portal/payments');
    expect(response?.status()).toBe(200);
    await page.screenshot({ path: '../docs/screenshots/portal-payments.png', fullPage: true });
  });
});
