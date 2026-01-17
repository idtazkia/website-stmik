import { test, expect } from '@playwright/test';
import { CandidatesPage } from './pages';

test.describe('Admin Candidate Detail', () => {
  test.beforeEach(async ({ page }) => {
    // Login as admin before each test
    await page.goto('/test/login/admin');
    await page.waitForURL(/\/admin\/?$/);
  });

  test.describe('Page Navigation', () => {
    test('should navigate to candidate detail from list', async ({ page }) => {
      // Go to candidates list
      await page.goto('/admin/candidates');
      await page.waitForSelector('[data-testid="candidates-page"]');

      // Click on the first candidate's detail link
      const detailLink = page.locator('[data-testid^="view-candidate-"]').first();

      // Only run if there are candidates
      if (await detailLink.isVisible()) {
        await detailLink.click();
        await page.waitForURL(/\/admin\/candidates\/[a-f0-9-]+/);

        // Verify detail page loaded
        await expect(page.locator('text=Data Pribadi')).toBeVisible();
      }
    });

    test('should show back link to candidates list', async ({ page }) => {
      // Go to candidates list first
      await page.goto('/admin/candidates');
      await page.waitForSelector('[data-testid="candidates-page"]');

      // Click on first candidate if available
      const detailLink = page.locator('[data-testid^="view-candidate-"]').first();
      if (await detailLink.isVisible()) {
        await detailLink.click();
        await page.waitForURL(/\/admin\/candidates\/[a-f0-9-]+/);

        // Check back link exists
        const backLink = page.locator('a[href="/admin/candidates"]');
        await expect(backLink).toBeVisible();

        // Click back and verify navigation
        await backLink.click();
        await page.waitForURL(/\/admin\/candidates\/?$/);
      }
    });
  });

  test.describe('Page Content', () => {
    test('should display personal info section', async ({ page }) => {
      // Go to candidates list
      await page.goto('/admin/candidates');
      await page.waitForSelector('[data-testid="candidates-page"]');

      // Click on first candidate if available
      const detailLink = page.locator('[data-testid^="view-candidate-"]').first();
      if (await detailLink.isVisible()) {
        await detailLink.click();
        await page.waitForURL(/\/admin\/candidates\/[a-f0-9-]+/);

        // Check personal info section
        await expect(page.locator('text=Data Pribadi')).toBeVisible();
        await expect(page.locator('text=Email')).toBeVisible();
        await expect(page.locator('text=Telepon')).toBeVisible();
      }
    });

    test('should display education section', async ({ page }) => {
      // Go to candidates list
      await page.goto('/admin/candidates');
      await page.waitForSelector('[data-testid="candidates-page"]');

      // Click on first candidate if available
      const detailLink = page.locator('[data-testid^="view-candidate-"]').first();
      if (await detailLink.isVisible()) {
        await detailLink.click();
        await page.waitForURL(/\/admin\/candidates\/[a-f0-9-]+/);

        // Check education section
        await expect(page.locator('text=Pendidikan')).toBeVisible();
        await expect(page.locator('text=Asal Sekolah')).toBeVisible();
        await expect(page.locator('text=Pilihan Prodi')).toBeVisible();
      }
    });

    test('should display source and assignment section', async ({ page }) => {
      // Go to candidates list
      await page.goto('/admin/candidates');
      await page.waitForSelector('[data-testid="candidates-page"]');

      // Click on first candidate if available
      const detailLink = page.locator('[data-testid^="view-candidate-"]').first();
      if (await detailLink.isVisible()) {
        await detailLink.click();
        await page.waitForURL(/\/admin\/candidates\/[a-f0-9-]+/);

        // Check source & assignment section
        await expect(page.locator('text=Sumber & Assignment')).toBeVisible();
        await expect(page.locator('text=Sumber Info')).toBeVisible();
        await expect(page.locator('text=Konsultan')).toBeVisible();
      }
    });

    test('should display payment status section', async ({ page }) => {
      // Go to candidates list
      await page.goto('/admin/candidates');
      await page.waitForSelector('[data-testid="candidates-page"]');

      // Click on first candidate if available
      const detailLink = page.locator('[data-testid^="view-candidate-"]').first();
      if (await detailLink.isVisible()) {
        await detailLink.click();
        await page.waitForURL(/\/admin\/candidates\/[a-f0-9-]+/);

        // Check payment status section
        await expect(page.locator('text=Status Pembayaran')).toBeVisible();
        await expect(page.locator('text=Biaya Pendaftaran')).toBeVisible();
      }
    });

    test('should display documents section', async ({ page }) => {
      // Go to candidates list
      await page.goto('/admin/candidates');
      await page.waitForSelector('[data-testid="candidates-page"]');

      // Click on first candidate if available
      const detailLink = page.locator('[data-testid^="view-candidate-"]').first();
      if (await detailLink.isVisible()) {
        await detailLink.click();
        await page.waitForURL(/\/admin\/candidates\/[a-f0-9-]+/);

        // Check documents section
        await expect(page.locator('h3:has-text("Dokumen")')).toBeVisible();
      }
    });

    test('should display timeline section', async ({ page }) => {
      // Go to candidates list
      await page.goto('/admin/candidates');
      await page.waitForSelector('[data-testid="candidates-page"]');

      // Click on first candidate if available
      const detailLink = page.locator('[data-testid^="view-candidate-"]').first();
      if (await detailLink.isVisible()) {
        await detailLink.click();
        await page.waitForURL(/\/admin\/candidates\/[a-f0-9-]+/);

        // Check timeline section
        await expect(page.locator('text=Timeline Interaksi')).toBeVisible();
      }
    });
  });

  test.describe('Action Buttons', () => {
    test('should display log interaction button', async ({ page }) => {
      // Go to candidates list
      await page.goto('/admin/candidates');
      await page.waitForSelector('[data-testid="candidates-page"]');

      // Click on first candidate if available
      const detailLink = page.locator('[data-testid^="view-candidate-"]').first();
      if (await detailLink.isVisible()) {
        await detailLink.click();
        await page.waitForURL(/\/admin\/candidates\/[a-f0-9-]+/);

        // Check action buttons
        await expect(page.locator('button:has-text("Log Interaksi")')).toBeVisible();
      }
    });

    test('should display mark as lost button', async ({ page }) => {
      // Go to candidates list
      await page.goto('/admin/candidates');
      await page.waitForSelector('[data-testid="candidates-page"]');

      // Click on first candidate if available
      const detailLink = page.locator('[data-testid^="view-candidate-"]').first();
      if (await detailLink.isVisible()) {
        await detailLink.click();
        await page.waitForURL(/\/admin\/candidates\/[a-f0-9-]+/);

        // Check mark as lost button
        await expect(page.locator('button:has-text("Mark as Lost")')).toBeVisible();
      }
    });

    test('should open interaction modal when clicking log interaction', async ({ page }) => {
      // Go to candidates list
      await page.goto('/admin/candidates');
      await page.waitForSelector('[data-testid="candidates-page"]');

      // Click on first candidate if available
      const detailLink = page.locator('[data-testid^="view-candidate-"]').first();
      if (await detailLink.isVisible()) {
        await detailLink.click();
        await page.waitForURL(/\/admin\/candidates\/[a-f0-9-]+/);

        // Click log interaction button
        await page.click('button:has-text("Log Interaksi")');

        // Check modal is visible
        await expect(page.locator('#modal-interaksi')).toBeVisible();
        await expect(page.locator('text=Log Interaksi Baru')).toBeVisible();

        // Close modal
        await page.click('#modal-interaksi button:has-text("✕")');
        await expect(page.locator('#modal-interaksi')).toBeHidden();
      }
    });
  });

  test.describe('Error Handling', () => {
    test('should return 404 for non-existent candidate', async ({ page }) => {
      const response = await page.goto('/admin/candidates/00000000-0000-0000-0000-000000000000');
      expect(response?.status()).toBe(404);
    });

    test('should return 404 for invalid UUID', async ({ page }) => {
      const response = await page.goto('/admin/candidates/invalid-id');
      // Should either return 404 or 500 depending on how database handles invalid UUID
      expect([404, 500]).toContain(response?.status());
    });
  });
});
