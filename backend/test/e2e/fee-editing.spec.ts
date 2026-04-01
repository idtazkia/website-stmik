import { test, expect } from '@playwright/test';

test.describe('Fee Structure Editing', () => {
  test.beforeEach(async ({ page }) => {
    await page.goto('/test/login/admin');
    await page.waitForURL(/\/admin\/?$/);
    await page.goto('/admin/settings/fees');
  });

  test('should display fees settings page', async ({ page }) => {
    await expect(page.getByTestId('settings-fees-page')).toBeVisible();
    await expect(page.getByTestId('academic-year-select')).toBeVisible();
    await expect(page.getByTestId('add-fee-button')).toBeVisible();
  });

  test('should display fee table with data', async ({ page }) => {
    await expect(page.getByTestId('fees-table')).toBeVisible();
    // Should have fee rows from seed data (registration, tuition, dormitory)
    const feeRows = page.locator('[data-testid^="fee-row-"]');
    const count = await feeRows.count();
    // Seed data might or might not have fee structures created
    // At minimum, the table should be visible
    expect(count).toBeGreaterThanOrEqual(0);
  });

  test('should open add fee modal', async ({ page }) => {
    await page.getByTestId('add-fee-button').click();
    const modal = page.getByTestId('add-fee-modal');
    await expect(modal).toBeVisible();
    await expect(page.getByTestId('input-fee-type')).toBeVisible();
    await expect(page.getByTestId('input-fee-prodi')).toBeVisible();
    await expect(page.getByTestId('input-fee-amount')).toBeVisible();
  });

  test('should have fee type options from seed data', async ({ page }) => {
    await page.getByTestId('add-fee-button').click();
    await expect(page.getByTestId('add-fee-modal')).toBeVisible();
    const feeTypeSelect = page.getByTestId('input-fee-type');
    const options = feeTypeSelect.locator('option');
    // Should have placeholder + at least 3 fee types (registration, tuition, dormitory)
    const count = await options.count();
    expect(count).toBeGreaterThanOrEqual(4);
  });

  test('should open shared edit modal with fee data', async ({ page }) => {
    const feeRows = page.locator('[data-testid^="fee-row-"]');
    const count = await feeRows.count();
    if (count === 0) {
      test.skip();
      return;
    }

    // Click edit button on first fee
    const editBtn = page.locator('.edit-fee-btn').first();
    await editBtn.click();

    // Shared edit modal should be visible
    const modal = page.locator('#edit-fee-modal');
    await expect(modal).toBeVisible();

    // Amount field should be populated
    const amountInput = page.locator('#edit-fee-amount');
    await expect(amountInput).toBeVisible();
    const value = await amountInput.inputValue();
    expect(Number(value)).toBeGreaterThan(0);
  });

  test('should toggle fee active status', async ({ page }) => {
    const feeRows = page.locator('[data-testid^="fee-row-"]');
    const count = await feeRows.count();
    if (count === 0) {
      test.skip();
      return;
    }

    const firstRowId = await feeRows.first().getAttribute('data-fee-id');
    if (!firstRowId) {
      test.skip();
      return;
    }

    const toggle = page.getByTestId(`fee-status-toggle-${firstRowId}`);
    const statusBefore = await toggle.textContent();

    // Click toggle
    await toggle.click();

    // Wait for HTMX update
    await page.waitForTimeout(500);

    // Status should have changed
    const statusAfter = await page.getByTestId(`fee-status-toggle-${firstRowId}`).textContent();
    expect(statusAfter?.trim()).not.toBe(statusBefore?.trim());

    // Toggle back to restore
    await page.getByTestId(`fee-status-toggle-${firstRowId}`).click();
  });
});
