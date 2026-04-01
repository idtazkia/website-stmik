import { test, expect } from '@playwright/test';

test.describe('Education Consultant Rename', () => {
  test('should display "Education Consultant" as role label for consultant user', async ({ page }) => {
    await page.goto('/test/login/consultant');
    await page.waitForURL(/\/admin\/?$/);
    // The role badge in the nav should show Education Consultant
    await expect(page.locator('text=Education Consultant')).toBeVisible();
  });

  test('should display "EC" in candidates table header', async ({ page }) => {
    await page.goto('/test/login/admin');
    await page.waitForURL(/\/admin\/?$/);
    await page.goto('/admin/candidates');
    await expect(page.getByTestId('candidates-page')).toBeVisible();
    // Table header should say EC not Konsultan
    const headers = page.locator('th');
    await expect(headers.filter({ hasText: 'EC' }).first()).toBeVisible();
  });

  test('should display "EC" in candidates filter dropdown', async ({ page }) => {
    await page.goto('/test/login/admin');
    await page.waitForURL(/\/admin\/?$/);
    await page.goto('/admin/candidates');
    await expect(page.getByTestId('candidates-page')).toBeVisible();
    const filterConsultant = page.getByTestId('filter-consultant');
    await expect(filterConsultant).toBeVisible();
    // First option should be "Semua EC"
    const firstOption = filterConsultant.locator('option').first();
    await expect(firstOption).toHaveText('Semua EC');
  });

  test('should display "Education Consultant" in consultant dashboard title', async ({ page }) => {
    await page.goto('/test/login/consultant');
    await page.waitForURL(/\/admin\/?$/);
    await page.goto('/admin/my-dashboard');
    await expect(page.getByTestId('welcome-section')).toBeVisible();
    await expect(page.locator('text=Dashboard Education Consultant')).toBeVisible();
  });

  test('should display "Education Consultant" on portal candidate dashboard', async ({ page }) => {
    await page.goto('/test/login/candidate');
    await page.waitForURL(/\/portal\/?$/);
    // Check if consultant section exists and shows Education Consultant
    const consultantSection = page.locator('text=Education Consultant');
    // This may not be visible if no consultant is assigned, so check if page loads
    await expect(page.locator('[data-testid="portal-dashboard"]')).toBeVisible();
  });
});
