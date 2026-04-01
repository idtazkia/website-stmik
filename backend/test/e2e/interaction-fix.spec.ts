import { test, expect } from '@playwright/test';

test.describe('Interaction Logging Fix', () => {
  test.beforeEach(async ({ page }) => {
    await page.goto('/test/login/consultant');
    await page.waitForURL(/\/admin\/?$/);
  });

  test('should navigate to dedicated interaction form from candidate detail', async ({ page }) => {
    await page.goto('/admin/candidates');
    await page.waitForSelector('[data-testid="candidates-page"]');

    const detailLink = page.locator('[data-testid^="view-candidate-"]').first();
    if (!(await detailLink.isVisible())) {
      test.skip();
      return;
    }

    await detailLink.click();
    await page.waitForURL(/\/admin\/candidates\/[a-f0-9-]+/);

    // Log Interaksi button should be a link (not a modal trigger)
    const logBtn = page.getByTestId('btn-log-interaction');
    await expect(logBtn).toBeVisible();

    // Should be an <a> tag linking to the interaction form page
    const tagName = await logBtn.evaluate(el => el.tagName.toLowerCase());
    expect(tagName).toBe('a');

    // Click and verify navigation to interaction form page
    await logBtn.click();
    await page.waitForURL(/\/admin\/candidates\/[a-f0-9-]+\/interaction/);
    await expect(page.getByTestId('interaction-form')).toBeVisible();
  });

  test('should NOT have broken modal on candidate detail page', async ({ page }) => {
    await page.goto('/admin/candidates');
    await page.waitForSelector('[data-testid="candidates-page"]');

    const detailLink = page.locator('[data-testid^="view-candidate-"]').first();
    if (!(await detailLink.isVisible())) {
      test.skip();
      return;
    }

    await detailLink.click();
    await page.waitForURL(/\/admin\/candidates\/[a-f0-9-]+/);

    // The old broken modal should NOT exist
    await expect(page.locator('#modal-interaksi')).toHaveCount(0);
    await expect(page.getByTestId('modal-interaction')).toHaveCount(0);
  });

  test('should submit interaction form and show in timeline', async ({ page }) => {
    await page.goto('/admin/candidates');
    await page.waitForSelector('[data-testid="candidates-page"]');

    const detailLink = page.locator('[data-testid^="view-candidate-"]').first();
    if (!(await detailLink.isVisible())) {
      test.skip();
      return;
    }

    await detailLink.click();
    await page.waitForURL(/\/admin\/candidates\/[a-f0-9-]+/);

    const url = page.url();
    const candidateId = url.split('/').pop();

    // Navigate to interaction form
    await page.goto(`/admin/candidates/${candidateId}/interaction`);
    await expect(page.getByTestId('interaction-form')).toBeVisible();

    // Fill form
    await page.click('label:has(input[name="channel"][value="whatsapp"])');
    await page.click('label:has(input[name="category"]):first-of-type');

    const uniqueRemarks = `Fix verification test ${Date.now()} - testing form submission works correctly.`;
    await page.fill('textarea[name="remarks"]', uniqueRemarks);

    // Submit
    await page.click('button[value="save"]');

    // Should redirect to candidate detail
    await page.waitForURL(new RegExp(`/admin/candidates/${candidateId}$`), { timeout: 10000 });

    // Verify timeline shows the interaction
    await expect(page.locator('text=Timeline Interaksi')).toBeVisible();
    await expect(page.locator(`text=${uniqueRemarks.substring(0, 30)}`)).toBeVisible();
  });

  test('should persist interaction after page reload', async ({ page }) => {
    await page.goto('/admin/candidates');
    await page.waitForSelector('[data-testid="candidates-page"]');

    const detailLink = page.locator('[data-testid^="view-candidate-"]').first();
    if (!(await detailLink.isVisible())) {
      test.skip();
      return;
    }

    await detailLink.click();
    await page.waitForURL(/\/admin\/candidates\/[a-f0-9-]+/);

    const url = page.url();
    const candidateId = url.split('/').pop();

    await page.goto(`/admin/candidates/${candidateId}/interaction`);
    await expect(page.getByTestId('interaction-form')).toBeVisible();

    await page.click('label:has(input[name="channel"][value="call"])');
    await page.click('label:has(input[name="category"]):first-of-type');

    const uniqueRemarks = `Persistence check ${Date.now()} - this should survive a reload.`;
    await page.fill('textarea[name="remarks"]', uniqueRemarks);

    await page.click('button[value="save"]');
    await page.waitForURL(new RegExp(`/admin/candidates/${candidateId}$`), { timeout: 10000 });

    // Reload
    await page.reload();

    // Should still be visible
    await expect(page.locator(`text=${uniqueRemarks.substring(0, 20)}`)).toBeVisible();
  });
});
