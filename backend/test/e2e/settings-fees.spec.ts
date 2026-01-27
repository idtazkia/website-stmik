import { test, expect } from '@playwright/test';
import { SettingsFeesPage } from './pages';

test.describe('Settings - Fee Structure Management', () => {
  let feesPage: SettingsFeesPage;

  test.beforeEach(async ({ page }) => {
    feesPage = new SettingsFeesPage(page);
    // Login as admin before each test
    await feesPage.login('admin');
    await feesPage.goto();
    await feesPage.expectPageLoaded();
  });

  test.describe('Page Load', () => {
    test('should display fees page with table', async () => {
      await expect(feesPage.feesSection).toBeVisible();
      await expect(feesPage.feesTable).toBeVisible();
    });

    test('should display academic year selector', async () => {
      await expect(feesPage.academicYearSelect).toBeVisible();
    });

    test('should display add fee button', async () => {
      await expect(feesPage.addFeeButton).toBeVisible();
    });
  });

  test.describe('Fee Display', () => {
    test('should display fee type columns in table header', async () => {
      // Verify table headers exist
      const table = feesPage.feesTable;
      await expect(table.locator('th').first()).toBeVisible();
    });
  });

  test.describe('Fee CRUD', () => {
    // Run CRUD tests serially to avoid race conditions
    test.describe.configure({ mode: 'serial' });

    test('should open add fee modal', async () => {
      await feesPage.openAddFeeModal();
      await expect(feesPage.addFeeModal).toBeVisible();
      await expect(feesPage.inputFeeType).toBeVisible();
      await expect(feesPage.inputFeeProdi).toBeVisible();
      await expect(feesPage.inputFeeAmount).toBeVisible();
    });

    test('should add new fee structure via HTMX', async ({ page }) => {
      // Get current fee count
      const feeIdsBefore = await feesPage.getAllFeeIds();
      const countBefore = feeIdsBefore.length;

      // Add new fee (first fee type, all prodi, 5000000)
      await feesPage.addFee(1, null, 5000000);

      // Verify new fee appears
      const feeIdsAfter = await feesPage.getAllFeeIds();
      expect(feeIdsAfter.length).toBe(countBefore + 1);

      // Find the new fee
      const newFeeId = feeIdsAfter.find(id => !feeIdsBefore.includes(id));
      expect(newFeeId).toBeTruthy();

      if (newFeeId) {
        await feesPage.expectFeeDisplayed(newFeeId);

        // Reload and verify persistence
        await page.reload();
        await feesPage.expectPageLoaded();
        await feesPage.expectFeeDisplayed(newFeeId);
      }
    });

    test('should toggle fee status via HTMX', async () => {
      const feeIds = await feesPage.getAllFeeIds();
      if (feeIds.length === 0) {
        test.skip();
        return;
      }

      const feeId = feeIds[0];

      // Get current status
      const statusBefore = await feesPage.getFeeStatusToggle(feeId).textContent();
      const isActiveBefore = statusBefore?.trim() === 'Aktif';

      // Toggle status
      await feesPage.toggleFeeStatus(feeId);

      // Verify status changed
      if (isActiveBefore) {
        await feesPage.expectFeeStatus(feeId, 'inactive');
      } else {
        await feesPage.expectFeeStatus(feeId, 'active');
      }

      // Toggle back to restore original state
      await feesPage.toggleFeeStatus(feeId);
    });

    test('should edit fee amount via HTMX', async ({ page }) => {
      const feeIds = await feesPage.getAllFeeIds();
      if (feeIds.length === 0) {
        test.skip();
        return;
      }

      const feeId = feeIds[0];

      // Edit with new amount
      const newAmount = 7500000;
      await feesPage.editFeeAmount(feeId, newAmount);

      // Verify amount changed (check for formatted amount)
      await feesPage.expectFeeAmountContains(feeId, '7.500.000');

      // Reload and verify persistence
      await page.reload();
      await feesPage.expectPageLoaded();
      await feesPage.expectFeeAmountContains(feeId, '7.500.000');
    });

    test('should display edit button for each fee', async () => {
      const feeIds = await feesPage.getAllFeeIds();
      if (feeIds.length === 0) {
        test.skip();
        return;
      }

      for (const feeId of feeIds) {
        await expect(feesPage.getFeeEditButton(feeId)).toBeVisible();
      }
    });
  });
});
