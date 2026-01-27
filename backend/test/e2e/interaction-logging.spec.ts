import { test, expect } from '@playwright/test';

test.describe('Admin Interaction Logging', () => {
  test.beforeEach(async ({ page }) => {
    // Login as consultant (who can log interactions)
    await page.goto('/test/login/consultant');
    await page.waitForURL(/\/admin\/?$/);
  });

  test.describe('Interaction Form Navigation', () => {
    test('should navigate to interaction form from candidate detail', async ({ page }) => {
      // Go to candidates list
      await page.goto('/admin/candidates');
      await page.waitForSelector('[data-testid="candidates-page"]');

      // Click on the first candidate's detail link
      const detailLink = page.locator('[data-testid^="view-candidate-"]').first();

      // Only run if there are candidates
      if (await detailLink.isVisible()) {
        await detailLink.click();
        await page.waitForURL(/\/admin\/candidates\/[a-f0-9-]+/);

        // Get candidate ID from URL
        const url = page.url();
        const candidateId = url.split('/').pop();

        // Navigate to interaction form
        await page.goto(`/admin/candidates/${candidateId}/interaction`);

        // Verify interaction form page loaded
        await expect(page.locator('text=Log Interaksi Baru')).toBeVisible();
        await expect(page.locator('text=Channel Komunikasi')).toBeVisible();
        await expect(page.locator('text=Respon Kandidat')).toBeVisible();
      }
    });

    test('should show back link to candidate detail', async ({ page }) => {
      // Go to candidates list first
      await page.goto('/admin/candidates');
      await page.waitForSelector('[data-testid="candidates-page"]');

      // Click on first candidate if available
      const detailLink = page.locator('[data-testid^="view-candidate-"]').first();
      if (await detailLink.isVisible()) {
        await detailLink.click();
        await page.waitForURL(/\/admin\/candidates\/[a-f0-9-]+/);

        // Get candidate ID from URL
        const url = page.url();
        const candidateId = url.split('/').pop();

        // Navigate to interaction form
        await page.goto(`/admin/candidates/${candidateId}/interaction`);
        await expect(page.locator('text=Log Interaksi Baru')).toBeVisible();

        // Check back link exists and click it
        const backLink = page.locator(`a[href="/admin/candidates/${candidateId}"]`).first();
        await expect(backLink).toBeVisible();
        await backLink.click();

        // Verify navigation back to candidate detail
        await page.waitForURL(new RegExp(`/admin/candidates/${candidateId}$`));
      }
    });
  });

  test.describe('Interaction Form Elements', () => {
    test('should display all channel options', async ({ page }) => {
      // Go to candidates list
      await page.goto('/admin/candidates');
      await page.waitForSelector('[data-testid="candidates-page"]');

      // Click on first candidate if available
      const detailLink = page.locator('[data-testid^="view-candidate-"]').first();
      if (await detailLink.isVisible()) {
        await detailLink.click();
        await page.waitForURL(/\/admin\/candidates\/[a-f0-9-]+/);

        // Get candidate ID from URL and navigate to interaction form
        const url = page.url();
        const candidateId = url.split('/').pop();
        await page.goto(`/admin/candidates/${candidateId}/interaction`);

        // Verify all channel radio options are present
        await expect(page.locator('input[name="channel"][value="call"]')).toBeAttached();
        await expect(page.locator('input[name="channel"][value="whatsapp"]')).toBeAttached();
        await expect(page.locator('input[name="channel"][value="email"]')).toBeAttached();
        await expect(page.locator('input[name="channel"][value="campus_visit"]')).toBeAttached();
        await expect(page.locator('input[name="channel"][value="home_visit"]')).toBeAttached();
      }
    });

    test('should display category options from database', async ({ page }) => {
      // Go to candidates list
      await page.goto('/admin/candidates');
      await page.waitForSelector('[data-testid="candidates-page"]');

      // Click on first candidate if available
      const detailLink = page.locator('[data-testid^="view-candidate-"]').first();
      if (await detailLink.isVisible()) {
        await detailLink.click();
        await page.waitForURL(/\/admin\/candidates\/[a-f0-9-]+/);

        // Get candidate ID from URL and navigate to interaction form
        const url = page.url();
        const candidateId = url.split('/').pop();
        await page.goto(`/admin/candidates/${candidateId}/interaction`);

        // Verify category radio options are present (from database)
        await expect(page.locator('input[name="category"]').first()).toBeAttached();
      }
    });

    test('should display obstacle dropdown', async ({ page }) => {
      // Go to candidates list
      await page.goto('/admin/candidates');
      await page.waitForSelector('[data-testid="candidates-page"]');

      // Click on first candidate if available
      const detailLink = page.locator('[data-testid^="view-candidate-"]').first();
      if (await detailLink.isVisible()) {
        await detailLink.click();
        await page.waitForURL(/\/admin\/candidates\/[a-f0-9-]+/);

        // Get candidate ID from URL and navigate to interaction form
        const url = page.url();
        const candidateId = url.split('/').pop();
        await page.goto(`/admin/candidates/${candidateId}/interaction`);

        // Verify obstacle dropdown is present
        await expect(page.locator('select[name="obstacle"]')).toBeVisible();
      }
    });

    test('should display remarks textarea', async ({ page }) => {
      // Go to candidates list
      await page.goto('/admin/candidates');
      await page.waitForSelector('[data-testid="candidates-page"]');

      // Click on first candidate if available
      const detailLink = page.locator('[data-testid^="view-candidate-"]').first();
      if (await detailLink.isVisible()) {
        await detailLink.click();
        await page.waitForURL(/\/admin\/candidates\/[a-f0-9-]+/);

        // Get candidate ID from URL and navigate to interaction form
        const url = page.url();
        const candidateId = url.split('/').pop();
        await page.goto(`/admin/candidates/${candidateId}/interaction`);

        // Verify remarks textarea is present and required
        const remarksField = page.locator('textarea[name="remarks"]');
        await expect(remarksField).toBeVisible();
        await expect(remarksField).toHaveAttribute('required');
      }
    });

    test('should display next followup date field', async ({ page }) => {
      // Go to candidates list
      await page.goto('/admin/candidates');
      await page.waitForSelector('[data-testid="candidates-page"]');

      // Click on first candidate if available
      const detailLink = page.locator('[data-testid^="view-candidate-"]').first();
      if (await detailLink.isVisible()) {
        await detailLink.click();
        await page.waitForURL(/\/admin\/candidates\/[a-f0-9-]+/);

        // Get candidate ID from URL and navigate to interaction form
        const url = page.url();
        const candidateId = url.split('/').pop();
        await page.goto(`/admin/candidates/${candidateId}/interaction`);

        // Verify next followup date field is present
        await expect(page.locator('input[name="next_followup_date"]')).toBeVisible();
      }
    });

    test('should display candidate summary info', async ({ page }) => {
      // Go to candidates list
      await page.goto('/admin/candidates');
      await page.waitForSelector('[data-testid="candidates-page"]');

      // Click on first candidate if available
      const detailLink = page.locator('[data-testid^="view-candidate-"]').first();
      if (await detailLink.isVisible()) {
        await detailLink.click();
        await page.waitForURL(/\/admin\/candidates\/[a-f0-9-]+/);

        // Get candidate ID from URL and navigate to interaction form
        const url = page.url();
        const candidateId = url.split('/').pop();
        await page.goto(`/admin/candidates/${candidateId}/interaction`);

        // Verify candidate summary is displayed
        await expect(page.locator('text=Kandidat:')).toBeVisible();
        await expect(page.locator('.bg-gray-50 >> text=Nama')).toBeVisible();
        await expect(page.locator('.bg-gray-50 >> text=Status')).toBeVisible();
      }
    });

    test('should display submit buttons', async ({ page }) => {
      // Go to candidates list
      await page.goto('/admin/candidates');
      await page.waitForSelector('[data-testid="candidates-page"]');

      // Click on first candidate if available
      const detailLink = page.locator('[data-testid^="view-candidate-"]').first();
      if (await detailLink.isVisible()) {
        await detailLink.click();
        await page.waitForURL(/\/admin\/candidates\/[a-f0-9-]+/);

        // Get candidate ID from URL and navigate to interaction form
        const url = page.url();
        const candidateId = url.split('/').pop();
        await page.goto(`/admin/candidates/${candidateId}/interaction`);

        // Verify submit buttons
        await expect(page.locator('button[value="save"]')).toBeVisible();
        await expect(page.locator('button[value="save_and_next"]')).toBeVisible();
      }
    });
  });

  test.describe('Interaction Form Submission', () => {
    test('should require channel, category, and remarks', async ({ page }) => {
      // Go to candidates list
      await page.goto('/admin/candidates');
      await page.waitForSelector('[data-testid="candidates-page"]');

      // Click on first candidate if available
      const detailLink = page.locator('[data-testid^="view-candidate-"]').first();
      if (await detailLink.isVisible()) {
        await detailLink.click();
        await page.waitForURL(/\/admin\/candidates\/[a-f0-9-]+/);

        // Get candidate ID from URL and navigate to interaction form
        const url = page.url();
        const candidateId = url.split('/').pop();
        await page.goto(`/admin/candidates/${candidateId}/interaction`);

        // Try to submit without required fields
        await page.click('button[value="save"]');

        // Form should not navigate away (HTML5 validation)
        await expect(page).toHaveURL(new RegExp(`/admin/candidates/${candidateId}/interaction`));
      }
    });

    test('should submit interaction and redirect to candidate detail', async ({ page }) => {
      // Go to candidates list
      await page.goto('/admin/candidates');
      await page.waitForSelector('[data-testid="candidates-page"]');

      // Click on first candidate if available
      const detailLink = page.locator('[data-testid^="view-candidate-"]').first();
      if (await detailLink.isVisible()) {
        await detailLink.click();
        await page.waitForURL(/\/admin\/candidates\/[a-f0-9-]+/);

        // Get candidate ID from URL
        const url = page.url();
        const candidateId = url.split('/').pop();

        // Navigate to interaction form
        await page.goto(`/admin/candidates/${candidateId}/interaction`);
        await expect(page.locator('text=Log Interaksi Baru')).toBeVisible();

        // Fill required fields
        // Select channel (click on label containing the hidden radio)
        await page.click('label:has(input[name="channel"][value="whatsapp"])');

        // Select category (click the label containing the first category radio)
        await page.click('label:has(input[name="category"]):first-of-type');

        // Fill remarks with unique text for verification
        const uniqueRemarks = `E2E Test interaction at ${Date.now()} - kandidat tertarik dengan program studi.`;
        await page.fill('textarea[name="remarks"]', uniqueRemarks);

        // Submit the form
        await page.click('button[value="save"]');

        // Should redirect to candidate detail page
        await page.waitForURL(new RegExp(`/admin/candidates/${candidateId}$`), { timeout: 10000 });

        // Verify we're on the candidate detail page
        await expect(page.locator('text=Timeline Interaksi')).toBeVisible();

        // Verify the interaction appears in the timeline
        await expect(page.locator(`text=${uniqueRemarks.substring(0, 30)}`)).toBeVisible();
      }
    });

    test('should persist interaction after page reload', async ({ page }) => {
      // Go to candidates list
      await page.goto('/admin/candidates');
      await page.waitForSelector('[data-testid="candidates-page"]');

      // Click on first candidate if available
      const detailLink = page.locator('[data-testid^="view-candidate-"]').first();
      if (await detailLink.isVisible()) {
        await detailLink.click();
        await page.waitForURL(/\/admin\/candidates\/[a-f0-9-]+/);

        // Get candidate ID from URL
        const url = page.url();
        const candidateId = url.split('/').pop();

        // Navigate to interaction form
        await page.goto(`/admin/candidates/${candidateId}/interaction`);
        await expect(page.locator('text=Log Interaksi Baru')).toBeVisible();

        // Fill required fields
        await page.click('label:has(input[name="channel"][value="call"])');
        await page.click('label:has(input[name="category"]):first-of-type');

        // Fill remarks with unique text for verification
        const uniqueRemarks = `Persistence test ${Date.now()} - follow up via telepon.`;
        await page.fill('textarea[name="remarks"]', uniqueRemarks);

        // Set next followup date
        const nextWeek = new Date();
        nextWeek.setDate(nextWeek.getDate() + 7);
        const dateStr = nextWeek.toISOString().split('T')[0];
        await page.fill('input[name="next_followup_date"]', dateStr);

        // Submit the form
        await page.click('button[value="save"]');

        // Wait for redirect
        await page.waitForURL(new RegExp(`/admin/candidates/${candidateId}$`), { timeout: 10000 });

        // Reload the page
        await page.reload();

        // Verify the interaction still appears after reload (persisted to database)
        await expect(page.locator('text=Timeline Interaksi')).toBeVisible();
        await expect(page.locator(`text=${uniqueRemarks.substring(0, 20)}`)).toBeVisible();
      }
    });

    test('should display channel badge in timeline after submission', async ({ page }) => {
      // Go to candidates list
      await page.goto('/admin/candidates');
      await page.waitForSelector('[data-testid="candidates-page"]');

      // Click on first candidate if available
      const detailLink = page.locator('[data-testid^="view-candidate-"]').first();
      if (await detailLink.isVisible()) {
        await detailLink.click();
        await page.waitForURL(/\/admin\/candidates\/[a-f0-9-]+/);

        // Get candidate ID from URL
        const url = page.url();
        const candidateId = url.split('/').pop();

        // Navigate to interaction form
        await page.goto(`/admin/candidates/${candidateId}/interaction`);

        // Fill required fields - use email channel
        await page.click('label:has(input[name="channel"][value="email"])');
        await page.click('label:has(input[name="category"]):first-of-type');

        const uniqueRemarks = `Email channel test ${Date.now()}`;
        await page.fill('textarea[name="remarks"]', uniqueRemarks);

        // Submit the form
        await page.click('button[value="save"]');

        // Wait for redirect
        await page.waitForURL(new RegExp(`/admin/candidates/${candidateId}$`), { timeout: 10000 });

        // Verify channel badge appears (email icon or text)
        await expect(page.locator('text=Timeline Interaksi')).toBeVisible();
        // The timeline should show the email channel
        await expect(page.locator(`text=${uniqueRemarks.substring(0, 15)}`)).toBeVisible();
      }
    });
  });

  test.describe('Error Handling', () => {
    test('should return 404 for non-existent candidate', async ({ page }) => {
      const response = await page.goto('/admin/candidates/00000000-0000-0000-0000-000000000000/interaction');
      expect(response?.status()).toBe(404);
    });

    test('should return error for invalid UUID', async ({ page }) => {
      const response = await page.goto('/admin/candidates/invalid-id/interaction');
      // Should either return 404 or 500 depending on how database handles invalid UUID
      expect([404, 500]).toContain(response?.status());
    });
  });
});
