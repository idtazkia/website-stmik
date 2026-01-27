import { test, expect } from '@playwright/test';

test.describe('Supervisor Dashboard', () => {
  test.beforeEach(async ({ page }) => {
    // Login as supervisor
    await page.goto('/test/login/supervisor');
  });

  test('supervisor can access dashboard', async ({ page }) => {
    await page.goto('/admin/supervisor-dashboard');

    // Check welcome section with "Dashboard Supervisor" subtitle
    await expect(page.getByTestId('supervisor-dashboard')).toBeVisible();
    await expect(page.getByTestId('page-title')).toContainText('Dashboard Supervisor');
  });

  test('supervisor dashboard shows team stats', async ({ page }) => {
    await page.goto('/admin/supervisor-dashboard');

    // Check for team stats section
    const statsSection = page.getByTestId('team-stats');
    await expect(statsSection).toBeVisible();

    // Check for stat labels within the stats section
    await expect(statsSection.getByText('Registered')).toBeVisible();
    await expect(statsSection.getByText('Prospecting')).toBeVisible();
    await expect(statsSection.getByText('Committed')).toBeVisible();
    await expect(statsSection.getByText('Enrolled')).toBeVisible();
    await expect(statsSection.getByText('Lost')).toBeVisible();
  });

  test('supervisor dashboard shows team performance table', async ({ page }) => {
    await page.goto('/admin/supervisor-dashboard');

    // Check for team performance section
    await expect(page.getByTestId('team-performance-section')).toBeVisible();

    // Check for table headers
    await expect(page.getByRole('columnheader', { name: /konsultan/i })).toBeVisible();
    await expect(page.getByRole('columnheader', { name: /aktif/i })).toBeVisible();
    await expect(page.getByRole('columnheader', { name: /hari ini/i })).toBeVisible();
    await expect(page.getByRole('columnheader', { name: /overdue/i })).toBeVisible();
  });

  test('supervisor dashboard shows stuck candidates section', async ({ page }) => {
    await page.goto('/admin/supervisor-dashboard');

    // Check for stuck candidates section
    await expect(page.getByTestId('stuck-candidates-section')).toBeVisible();

    // Check for section header
    await expect(page.getByRole('heading', { name: /kandidat stuck/i, level: 3 })).toBeVisible();
  });

  test('supervisor dashboard shows monthly performance', async ({ page }) => {
    await page.goto('/admin/supervisor-dashboard');

    // Check for monthly performance section
    await expect(page.getByTestId('monthly-performance-section')).toBeVisible();

    // Check for monthly stats labels
    await expect(page.getByText(/leads baru/i)).toBeVisible();
    await expect(page.getByText(/enrollments/i)).toBeVisible();
  });

  test('supervisor dashboard shows quick actions', async ({ page }) => {
    await page.goto('/admin/supervisor-dashboard');

    // Check for quick actions section
    await expect(page.getByTestId('quick-actions-section')).toBeVisible();

    // Check for quick action links
    await expect(page.getByRole('link', { name: /lihat kandidat stuck/i })).toBeVisible();
    await expect(page.getByRole('link', { name: /lihat report konsultan/i })).toBeVisible();
    await expect(page.getByRole('link', { name: /lihat funnel report/i })).toBeVisible();
  });

  test('supervisor sees team dashboard link in sidebar', async ({ page }) => {
    await page.goto('/admin');

    // Check for supervisor dashboard link in sidebar
    await expect(page.getByTestId('nav-supervisor-dashboard')).toBeVisible();
    await expect(page.getByText('Dashboard Tim')).toBeVisible();
  });

  test('consultant does not see team dashboard link', async ({ page }) => {
    // Login as consultant instead
    await page.goto('/test/login/consultant');
    await page.goto('/admin');

    // Consultant should not see supervisor dashboard link
    await expect(page.getByTestId('nav-supervisor-dashboard')).not.toBeVisible();
  });
});
