import { test, expect } from '@playwright/test';

test.describe('Bulk EC Assignment', () => {
  test.beforeEach(async ({ page }) => {
    await page.goto('/test/login/admin');
    await page.waitForURL(/\/admin\/?$/);
    await page.goto('/admin/candidates');
    await expect(page.getByTestId('candidates-page')).toBeVisible();
  });

  test('should display checkboxes for admin users', async ({ page }) => {
    // Admin should see select-all checkbox
    await expect(page.locator('#select-all')).toBeVisible();
    // Each candidate row should have a checkbox
    const checkboxes = page.locator('.candidate-checkbox');
    const count = await checkboxes.count();
    if (count > 0) {
      await expect(checkboxes.first()).toBeVisible();
    }
  });

  test('should show bulk action bar when candidates are selected', async ({ page }) => {
    const checkboxes = page.locator('.candidate-checkbox');
    const count = await checkboxes.count();
    if (count === 0) {
      test.skip();
      return;
    }

    // Initially bulk action bar should be hidden
    await expect(page.locator('#bulk-action-bar')).toBeHidden();

    // Select first candidate
    await checkboxes.first().check();

    // Bulk action bar should appear
    await expect(page.locator('#bulk-action-bar')).toBeVisible();
    await expect(page.locator('#bulk-count')).toContainText('1 kandidat dipilih');
  });

  test('should select all candidates with select-all checkbox', async ({ page }) => {
    const checkboxes = page.locator('.candidate-checkbox');
    const count = await checkboxes.count();
    if (count === 0) {
      test.skip();
      return;
    }

    await page.locator('#select-all').check();

    // All checkboxes should be checked
    for (let i = 0; i < count; i++) {
      await expect(checkboxes.nth(i)).toBeChecked();
    }

    // Bulk action bar should show correct count
    await expect(page.locator('#bulk-action-bar')).toBeVisible();
    await expect(page.locator('#bulk-count')).toContainText(`${count} kandidat dipilih`);
  });

  test('should clear selection with cancel button', async ({ page }) => {
    const checkboxes = page.locator('.candidate-checkbox');
    const count = await checkboxes.count();
    if (count === 0) {
      test.skip();
      return;
    }

    // Select all
    await page.locator('#select-all').check();
    await expect(page.locator('#bulk-action-bar')).toBeVisible();

    // Click cancel
    await page.locator('#bulk-cancel-btn').click();

    // Bar should be hidden, all unchecked
    await expect(page.locator('#bulk-action-bar')).toBeHidden();
    for (let i = 0; i < count; i++) {
      await expect(checkboxes.nth(i)).not.toBeChecked();
    }
  });

  test('should display EC options in bulk assign dropdown', async ({ page }) => {
    const select = page.locator('#bulk-consultant-select');
    await expect(select).toBeVisible();
    // Should have at least the placeholder option
    const options = select.locator('option');
    const optionCount = await options.count();
    expect(optionCount).toBeGreaterThanOrEqual(1);
  });

  test('should disable assign button when no EC selected', async ({ page }) => {
    const checkboxes = page.locator('.candidate-checkbox');
    const count = await checkboxes.count();
    if (count === 0) {
      test.skip();
      return;
    }

    await checkboxes.first().check();
    await expect(page.locator('#bulk-action-bar')).toBeVisible();

    // Assign button should be disabled when no EC selected
    await expect(page.locator('#bulk-assign-btn')).toBeDisabled();
  });

  test('should display "Belum Ditugaskan" filter option', async ({ page }) => {
    const filterConsultant = page.getByTestId('filter-consultant');
    await expect(filterConsultant).toBeVisible();
    const unassignedOption = filterConsultant.locator('option[value="unassigned"]');
    await expect(unassignedOption).toHaveText('Belum Ditugaskan');
  });

  test('consultant should NOT see checkboxes', async ({ page }) => {
    // Login as consultant instead
    await page.goto('/test/login/consultant');
    await page.waitForURL(/\/admin\/?$/);
    await page.goto('/admin/candidates');
    await expect(page.getByTestId('candidates-page')).toBeVisible();

    // Should NOT have select-all checkbox
    await expect(page.locator('#select-all')).toHaveCount(0);
  });
});
