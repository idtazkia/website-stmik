import { test, expect } from '@playwright/test';

// Screenshot tests for user manual documentation
// Run with: npx playwright test e2e/screenshots.spec.ts
// Screenshots saved to: docs/screenshots/

const SCREENSHOT_DIR = 'docs/screenshots';

test.describe('User Manual Screenshots - Admin Settings', () => {
  test.beforeEach(async ({ page }) => {
    await page.setViewportSize({ width: 1280, height: 900 });
    // Login as admin
    await page.goto('/test/login/admin');
    await page.waitForURL(/\/admin\/?$/);
  });

  // Authentication
  test('01 - Login Page', async ({ page }) => {
    await page.context().clearCookies();
    await page.goto('/admin/login');
    await expect(page.locator('text=Masuk dengan Google')).toBeVisible();
    await page.screenshot({ path: `${SCREENSHOT_DIR}/01-login.png`, fullPage: true });
  });

  // Dashboard (mockup - no specific test ID)
  test('02 - Admin Dashboard', async ({ page }) => {
    await page.goto('/admin');
    await expect(page.locator('text=Total Kandidat')).toBeVisible();
    await page.screenshot({ path: `${SCREENSHOT_DIR}/02-dashboard.png`, fullPage: true });
  });

  // Settings - Users
  test('03 - Settings Users', async ({ page }) => {
    await page.goto('/admin/settings/users');
    await expect(page.getByTestId('settings-users-page')).toBeVisible();
    await page.screenshot({ path: `${SCREENSHOT_DIR}/03-settings-users.png`, fullPage: true });
  });

  // Settings - Programs
  test('04 - Settings Programs', async ({ page }) => {
    await page.goto('/admin/settings/programs');
    await expect(page.getByTestId('settings-programs-page')).toBeVisible();
    await page.screenshot({ path: `${SCREENSHOT_DIR}/04-settings-programs.png`, fullPage: true });
  });

  test('04a - Settings Programs - Add Modal', async ({ page }) => {
    await page.goto('/admin/settings/programs');
    await expect(page.getByTestId('settings-programs-page')).toBeVisible();
    await page.getByTestId('add-program-button').click();
    await expect(page.getByTestId('add-program-modal')).toBeVisible();
    await page.screenshot({ path: `${SCREENSHOT_DIR}/04a-settings-programs-add.png`, fullPage: true });
  });

  // Settings - Fees
  test('05 - Settings Fees', async ({ page }) => {
    await page.goto('/admin/settings/fees');
    await expect(page.getByTestId('settings-fees-page')).toBeVisible();
    await page.screenshot({ path: `${SCREENSHOT_DIR}/05-settings-fees.png`, fullPage: true });
  });

  test('05a - Settings Fees - Add Modal', async ({ page }) => {
    await page.goto('/admin/settings/fees');
    await expect(page.getByTestId('settings-fees-page')).toBeVisible();
    await page.getByTestId('add-fee-button').click();
    await expect(page.getByTestId('add-fee-modal')).toBeVisible();
    await page.screenshot({ path: `${SCREENSHOT_DIR}/05a-settings-fees-add.png`, fullPage: true });
  });

  // Settings - Categories & Obstacles
  test('06 - Settings Categories', async ({ page }) => {
    await page.goto('/admin/settings/categories');
    await expect(page.getByTestId('settings-categories-page')).toBeVisible();
    await page.screenshot({ path: `${SCREENSHOT_DIR}/06-settings-categories.png`, fullPage: true });
  });

  test('06a - Settings Categories - Add Category Modal', async ({ page }) => {
    await page.goto('/admin/settings/categories');
    await expect(page.getByTestId('settings-categories-page')).toBeVisible();
    await page.getByTestId('add-category-button').click();
    await expect(page.getByTestId('add-category-modal')).toBeVisible();
    await page.screenshot({ path: `${SCREENSHOT_DIR}/06a-settings-categories-add.png`, fullPage: true });
  });

  test('06b - Settings Categories - Add Obstacle Modal', async ({ page }) => {
    await page.goto('/admin/settings/categories');
    await expect(page.getByTestId('settings-categories-page')).toBeVisible();
    await page.getByTestId('add-obstacle-button').click();
    await expect(page.getByTestId('add-obstacle-modal')).toBeVisible();
    await page.screenshot({ path: `${SCREENSHOT_DIR}/06b-settings-obstacles-add.png`, fullPage: true });
  });

  // Settings - Campaigns
  test('07 - Settings Campaigns', async ({ page }) => {
    await page.goto('/admin/settings/campaigns');
    await expect(page.getByTestId('settings-campaigns-page')).toBeVisible();
    await page.screenshot({ path: `${SCREENSHOT_DIR}/07-settings-campaigns.png`, fullPage: true });
  });

  test('07a - Settings Campaigns - Add Modal', async ({ page }) => {
    await page.goto('/admin/settings/campaigns');
    await expect(page.getByTestId('settings-campaigns-page')).toBeVisible();
    await page.getByTestId('add-campaign-button').click();
    await expect(page.getByTestId('add-campaign-modal')).toBeVisible();
    await page.screenshot({ path: `${SCREENSHOT_DIR}/07a-settings-campaigns-add.png`, fullPage: true });
  });

  // Settings - Rewards
  test('08 - Settings Rewards', async ({ page }) => {
    await page.goto('/admin/settings/rewards');
    await expect(page.getByTestId('settings-rewards-page')).toBeVisible();
    await page.screenshot({ path: `${SCREENSHOT_DIR}/08-settings-rewards.png`, fullPage: true });
  });

  test('08a - Settings Rewards - Add Reward Modal', async ({ page }) => {
    await page.goto('/admin/settings/rewards');
    await expect(page.getByTestId('settings-rewards-page')).toBeVisible();
    await page.getByTestId('add-reward-button').click();
    await expect(page.getByTestId('add-reward-modal')).toBeVisible();
    await page.screenshot({ path: `${SCREENSHOT_DIR}/08a-settings-rewards-add.png`, fullPage: true });
  });

  test('08b - Settings Rewards - Add MGM Modal', async ({ page }) => {
    await page.goto('/admin/settings/rewards');
    await expect(page.getByTestId('settings-rewards-page')).toBeVisible();
    await page.getByTestId('add-mgm-reward-button').click();
    await expect(page.getByTestId('add-mgm-reward-modal')).toBeVisible();
    await page.screenshot({ path: `${SCREENSHOT_DIR}/08b-settings-mgm-add.png`, fullPage: true });
  });

  // Settings - Referrers
  test('09 - Settings Referrers', async ({ page }) => {
    await page.goto('/admin/settings/referrers');
    await expect(page.getByTestId('settings-referrers-page')).toBeVisible();
    await page.screenshot({ path: `${SCREENSHOT_DIR}/09-settings-referrers.png`, fullPage: true });
  });

  test('09a - Settings Referrers - Add Modal', async ({ page }) => {
    await page.goto('/admin/settings/referrers');
    await expect(page.getByTestId('settings-referrers-page')).toBeVisible();
    await page.getByTestId('add-referrer-button').click();
    await expect(page.getByTestId('add-referrer-modal')).toBeVisible();
    await page.screenshot({ path: `${SCREENSHOT_DIR}/09a-settings-referrers-add.png`, fullPage: true });
  });

  // Settings - Assignment Algorithm
  test('10 - Settings Assignment', async ({ page }) => {
    await page.goto('/admin/settings/assignment');
    await expect(page.getByTestId('settings-assignment-page')).toBeVisible();
    await page.screenshot({ path: `${SCREENSHOT_DIR}/10-settings-assignment.png`, fullPage: true });
  });

  // Settings - Document Types
  test('11 - Settings Document Types', async ({ page }) => {
    await page.goto('/admin/settings/document-types');
    await expect(page.getByTestId('settings-documents-page')).toBeVisible();
    await page.screenshot({ path: `${SCREENSHOT_DIR}/11-settings-document-types.png`, fullPage: true });
  });

  test('11a - Settings Document Types - Add Modal', async ({ page }) => {
    await page.goto('/admin/settings/document-types');
    await expect(page.getByTestId('settings-documents-page')).toBeVisible();
    await page.getByTestId('add-document-type-button').click();
    await expect(page.getByTestId('add-document-type-modal')).toBeVisible();
    await page.screenshot({ path: `${SCREENSHOT_DIR}/11a-settings-document-types-add.png`, fullPage: true });
  });

  // Settings - Lost Reasons
  test('12 - Settings Lost Reasons', async ({ page }) => {
    await page.goto('/admin/settings/lost-reasons');
    await expect(page.getByTestId('settings-lost-reasons-page')).toBeVisible();
    await page.screenshot({ path: `${SCREENSHOT_DIR}/12-settings-lost-reasons.png`, fullPage: true });
  });

  test('12a - Settings Lost Reasons - Add Modal', async ({ page }) => {
    await page.goto('/admin/settings/lost-reasons');
    await expect(page.getByTestId('settings-lost-reasons-page')).toBeVisible();
    await page.getByTestId('add-lost-reason-button').click();
    await expect(page.getByTestId('add-lost-reason-modal')).toBeVisible();
    await page.screenshot({ path: `${SCREENSHOT_DIR}/12a-settings-lost-reasons-add.png`, fullPage: true });
  });
});
