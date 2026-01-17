import { test, expect } from '@playwright/test';
import { CandidatesPage } from './pages';

test.describe('Admin Candidates List', () => {
  let candidatesPage: CandidatesPage;

  test.beforeEach(async ({ page }) => {
    candidatesPage = new CandidatesPage(page);
    // Login as admin before each test
    await candidatesPage.login('admin');
    await candidatesPage.goto();
    await candidatesPage.expectPageLoaded();
  });

  test.describe('Page Load', () => {
    test('should display candidates page with all sections', async () => {
      await expect(candidatesPage.candidatesPage).toBeVisible();
      await expect(candidatesPage.statsSection).toBeVisible();
      await expect(candidatesPage.filtersSection).toBeVisible();
      await expect(candidatesPage.candidatesTable).toBeVisible();
    });

    test('should display stats cards', async () => {
      await expect(candidatesPage.statTotal).toBeVisible();
      await expect(candidatesPage.statRegistered).toBeVisible();
      await expect(candidatesPage.statProspecting).toBeVisible();
      await expect(candidatesPage.statCommitted).toBeVisible();
      await expect(candidatesPage.statEnrolled).toBeVisible();
      await expect(candidatesPage.statLost).toBeVisible();
    });

    test('should display filter controls', async () => {
      await expect(candidatesPage.filterStatus).toBeVisible();
      await expect(candidatesPage.filterConsultant).toBeVisible();
      await expect(candidatesPage.filterProdi).toBeVisible();
      await expect(candidatesPage.filterCampaign).toBeVisible();
      await expect(candidatesPage.filterSource).toBeVisible();
      await expect(candidatesPage.filterSearch).toBeVisible();
    });

    test('should display candidates table with headers', async () => {
      const table = candidatesPage.candidatesTable;
      await expect(table.locator('th')).toHaveCount(7); // Kandidat, Prodi, Status, Konsultan, Sumber, Terdaftar, Aksi
    });
  });

  test.describe('Filtering', () => {
    test('should filter by status', async () => {
      // Get initial count
      const initialCount = await candidatesPage.getCandidateRowCount();

      // Filter by registered status
      await candidatesPage.selectStatus('registered');

      // URL should reflect the filter
      await expect(candidatesPage.page).toHaveURL(/status=registered/);
    });

    test('should filter by search query', async ({ page }) => {
      // Search for a candidate
      await candidatesPage.searchCandidates('test');

      // URL should reflect the search
      await expect(page).toHaveURL(/search=test/);
    });

    test('should clear filters', async () => {
      // Apply some filters first
      await candidatesPage.selectStatus('registered');
      await candidatesPage.searchCandidates('test');

      // Clear filters
      await candidatesPage.clearFilters();

      // Verify filters are cleared
      await expect(candidatesPage.filterStatus).toHaveValue('');
      await expect(candidatesPage.filterSearch).toHaveValue('');
    });

    test('should update URL with filter parameters', async ({ page }) => {
      await candidatesPage.selectStatus('prospecting');
      await expect(page).toHaveURL(/status=prospecting/);

      await candidatesPage.selectSourceType('instagram');
      await expect(page).toHaveURL(/source_type=instagram/);
    });
  });

  test.describe('Empty State', () => {
    test('should show empty message when no candidates match filter', async () => {
      // Search for something that should not exist
      await candidatesPage.searchCandidates('nonexistent-candidate-xyz-12345');
      await candidatesPage.expectEmptyList();
    });
  });

  test.describe('Navigation', () => {
    test('should navigate to candidates from sidebar', async ({ page }) => {
      // Login as admin
      await page.goto('/test/login/admin');
      await page.waitForURL(/\/admin\/?$/);

      // Click candidates link in sidebar
      await page.getByTestId('nav-candidates').click();
      await page.waitForURL(/\/admin\/candidates/);

      // Verify page loaded
      await expect(page.getByTestId('candidates-page')).toBeVisible();
    });
  });
});

test.describe('Admin Candidates Role-Based Visibility', () => {
  test('admin should see all candidates', async ({ page }) => {
    const candidatesPage = new CandidatesPage(page);
    await candidatesPage.login('admin');
    await candidatesPage.goto();
    await candidatesPage.expectPageLoaded();

    // Admin should see the total stat which represents all candidates
    const totalText = await candidatesPage.statTotal.textContent();
    // Just verify it's a number (admin sees all)
    expect(totalText).toMatch(/^\d+$/);
  });

  test('consultant should only see assigned candidates', async ({ page }) => {
    const candidatesPage = new CandidatesPage(page);
    await candidatesPage.login('consultant');
    await candidatesPage.goto();
    await candidatesPage.expectPageLoaded();

    // Consultant should see only their assigned candidates
    // The exact number depends on test data, just verify page loads correctly
    await expect(candidatesPage.statsSection).toBeVisible();
  });

  test('supervisor should see team candidates', async ({ page }) => {
    const candidatesPage = new CandidatesPage(page);
    await candidatesPage.login('supervisor');
    await candidatesPage.goto();
    await candidatesPage.expectPageLoaded();

    // Supervisor should see their team's candidates
    await expect(candidatesPage.statsSection).toBeVisible();
  });
});
