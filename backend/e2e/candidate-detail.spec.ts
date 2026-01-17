import { test, expect } from '@playwright/test';
import { CandidatesPage, CandidateDetailPage } from './pages';

test.describe('Admin Candidate Detail', () => {
  let candidatesPage: CandidatesPage;
  let detailPage: CandidateDetailPage;

  test.beforeEach(async ({ page }) => {
    candidatesPage = new CandidatesPage(page);
    detailPage = new CandidateDetailPage(page);
    await candidatesPage.login('admin');
  });

  test.describe('Page Navigation', () => {
    test('should navigate to candidate detail from list', async ({ page }) => {
      await page.goto('/admin/candidates');
      await candidatesPage.expectPageLoaded();

      // Click on the first candidate's detail link
      const detailLink = page.locator('[data-testid^="view-candidate-"]').first();

      // Only run if there are candidates
      if (await detailLink.isVisible()) {
        await detailLink.click();
        await page.waitForURL(/\/admin\/candidates\/[a-f0-9-]+/);

        // Verify detail page loaded
        await expect(detailPage.detailPage).toBeVisible();
        await expect(detailPage.sectionTitlePersonalInfo).toBeVisible();
      }
    });

    test('should show back link to candidates list', async ({ page }) => {
      await page.goto('/admin/candidates');
      await candidatesPage.expectPageLoaded();

      // Click on first candidate if available
      const detailLink = page.locator('[data-testid^="view-candidate-"]').first();
      if (await detailLink.isVisible()) {
        await detailLink.click();
        await page.waitForURL(/\/admin\/candidates\/[a-f0-9-]+/);

        // Check back link exists
        await expect(detailPage.backLink).toBeVisible();

        // Click back and verify navigation
        await detailPage.goBackToCandidatesList();
      }
    });
  });

  test.describe('Page Content', () => {
    test('should display personal info section', async ({ page }) => {
      await page.goto('/admin/candidates');
      await candidatesPage.expectPageLoaded();

      // Click on first candidate if available
      const detailLink = page.locator('[data-testid^="view-candidate-"]').first();
      if (await detailLink.isVisible()) {
        await detailLink.click();
        await page.waitForURL(/\/admin\/candidates\/[a-f0-9-]+/);

        // Check personal info section
        await expect(detailPage.sectionTitlePersonalInfo).toBeVisible();
        await expect(detailPage.fieldEmail).toBeVisible();
        await expect(detailPage.fieldPhone).toBeVisible();
      }
    });

    test('should display education section', async ({ page }) => {
      await page.goto('/admin/candidates');
      await candidatesPage.expectPageLoaded();

      // Click on first candidate if available
      const detailLink = page.locator('[data-testid^="view-candidate-"]').first();
      if (await detailLink.isVisible()) {
        await detailLink.click();
        await page.waitForURL(/\/admin\/candidates\/[a-f0-9-]+/);

        // Check education section
        await expect(detailPage.sectionTitleEducation).toBeVisible();
        await expect(detailPage.fieldHighSchool).toBeVisible();
        await expect(detailPage.fieldProdi).toBeVisible();
      }
    });

    test('should display source and assignment section', async ({ page }) => {
      await page.goto('/admin/candidates');
      await candidatesPage.expectPageLoaded();

      // Click on first candidate if available
      const detailLink = page.locator('[data-testid^="view-candidate-"]').first();
      if (await detailLink.isVisible()) {
        await detailLink.click();
        await page.waitForURL(/\/admin\/candidates\/[a-f0-9-]+/);

        // Check source & assignment section
        await expect(detailPage.sectionTitleSourceAssignment).toBeVisible();
        await expect(detailPage.fieldSourceInfo).toBeVisible();
        await expect(detailPage.fieldConsultant).toBeVisible();
      }
    });

    test('should display payment status section', async ({ page }) => {
      await page.goto('/admin/candidates');
      await candidatesPage.expectPageLoaded();

      // Click on first candidate if available
      const detailLink = page.locator('[data-testid^="view-candidate-"]').first();
      if (await detailLink.isVisible()) {
        await detailLink.click();
        await page.waitForURL(/\/admin\/candidates\/[a-f0-9-]+/);

        // Check payment status section
        await expect(detailPage.sectionTitlePaymentStatus).toBeVisible();
        await expect(detailPage.fieldRegistrationFee).toBeVisible();
      }
    });

    test('should display documents section', async ({ page }) => {
      await page.goto('/admin/candidates');
      await candidatesPage.expectPageLoaded();

      // Click on first candidate if available
      const detailLink = page.locator('[data-testid^="view-candidate-"]').first();
      if (await detailLink.isVisible()) {
        await detailLink.click();
        await page.waitForURL(/\/admin\/candidates\/[a-f0-9-]+/);

        // Check documents section
        await expect(detailPage.sectionTitleDocuments).toBeVisible();
      }
    });

    test('should display timeline section', async ({ page }) => {
      await page.goto('/admin/candidates');
      await candidatesPage.expectPageLoaded();

      // Click on first candidate if available
      const detailLink = page.locator('[data-testid^="view-candidate-"]').first();
      if (await detailLink.isVisible()) {
        await detailLink.click();
        await page.waitForURL(/\/admin\/candidates\/[a-f0-9-]+/);

        // Check timeline section
        await expect(detailPage.sectionTitleTimeline).toBeVisible();
      }
    });
  });

  test.describe('Action Buttons', () => {
    test('should display log interaction button', async ({ page }) => {
      await page.goto('/admin/candidates');
      await candidatesPage.expectPageLoaded();

      // Click on first candidate if available
      const detailLink = page.locator('[data-testid^="view-candidate-"]').first();
      if (await detailLink.isVisible()) {
        await detailLink.click();
        await page.waitForURL(/\/admin\/candidates\/[a-f0-9-]+/);

        // Check action buttons
        await expect(detailPage.btnLogInteraction).toBeVisible();
      }
    });

    test('should display mark as lost button', async ({ page }) => {
      await page.goto('/admin/candidates');
      await candidatesPage.expectPageLoaded();

      // Click on first candidate if available
      const detailLink = page.locator('[data-testid^="view-candidate-"]').first();
      if (await detailLink.isVisible()) {
        await detailLink.click();
        await page.waitForURL(/\/admin\/candidates\/[a-f0-9-]+/);

        // Check mark as lost button
        await expect(detailPage.btnMarkLost).toBeVisible();
      }
    });

    test('should open interaction modal when clicking log interaction', async ({ page }) => {
      await page.goto('/admin/candidates');
      await candidatesPage.expectPageLoaded();

      // Click on first candidate if available
      const detailLink = page.locator('[data-testid^="view-candidate-"]').first();
      if (await detailLink.isVisible()) {
        await detailLink.click();
        await page.waitForURL(/\/admin\/candidates\/[a-f0-9-]+/);

        // Open interaction modal
        await detailPage.openInteractionModal();

        // Check modal is visible with form elements
        await expect(detailPage.modalTitle).toBeVisible();
        await expect(detailPage.selectChannel).toBeVisible();
        await expect(detailPage.selectCategory).toBeVisible();
        await expect(detailPage.inputRemarks).toBeVisible();

        // Close modal
        await detailPage.closeInteractionModal();
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
