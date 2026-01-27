import { test, expect } from '@playwright/test';

test.describe('Report Pages', () => {
  test.describe('Funnel Report', () => {
    test('admin can access funnel report', async ({ browser }) => {
      const adminPage = await browser.newPage();
      await adminPage.goto('/test/login/admin');
      await adminPage.goto('/admin/reports/funnel');

      await expect(adminPage.getByRole('heading', { name: 'Laporan Funnel', level: 2 })).toBeVisible();

      await adminPage.close();
    });

    test('funnel report shows stage data', async ({ browser }) => {
      const adminPage = await browser.newPage();
      await adminPage.goto('/test/login/admin');
      await adminPage.goto('/admin/reports/funnel');

      // Check that funnel visualization is displayed
      await expect(adminPage.locator('.flex.items-end.justify-center').first()).toBeVisible();

      await adminPage.close();
    });

    test('funnel report shows conversion rates', async ({ browser }) => {
      const adminPage = await browser.newPage();
      await adminPage.goto('/test/login/admin');
      await adminPage.goto('/admin/reports/funnel');

      // Check that conversion cards are displayed
      await expect(adminPage.locator('.grid.grid-cols-3')).toBeVisible();

      await adminPage.close();
    });
  });

  test.describe('Campaign Report', () => {
    test('admin can access campaign report', async ({ browser }) => {
      const adminPage = await browser.newPage();
      await adminPage.goto('/test/login/admin');
      await adminPage.goto('/admin/reports/campaigns');

      await expect(adminPage.getByRole('heading', { name: 'Laporan ROI Kampanye', level: 2 })).toBeVisible();

      await adminPage.close();
    });

    test('campaign report shows summary stats', async ({ browser }) => {
      const adminPage = await browser.newPage();
      await adminPage.goto('/test/login/admin');
      await adminPage.goto('/admin/reports/campaigns');

      // Check summary stats
      await expect(adminPage.getByText('Total Leads')).toBeVisible();
      await expect(adminPage.getByText('Total Enrolled')).toBeVisible();
      await expect(adminPage.getByText('Avg Conversion')).toBeVisible();

      await adminPage.close();
    });

    test('campaign report shows campaign table', async ({ browser }) => {
      const adminPage = await browser.newPage();
      await adminPage.goto('/test/login/admin');
      await adminPage.goto('/admin/reports/campaigns');

      // Check table headers using role
      await expect(adminPage.getByRole('columnheader', { name: 'Kampanye' })).toBeVisible();
      await expect(adminPage.getByRole('columnheader', { name: 'Tipe' })).toBeVisible();
      await expect(adminPage.getByRole('columnheader', { name: 'Conversion' })).toBeVisible();

      await adminPage.close();
    });
  });

  test.describe('Consultant Report', () => {
    test('admin can access consultant report', async ({ browser }) => {
      const adminPage = await browser.newPage();
      await adminPage.goto('/test/login/admin');
      await adminPage.goto('/admin/reports/consultants');

      await expect(adminPage.getByTestId('consultant-report-page')).toBeVisible();

      await adminPage.close();
    });

    test('consultant report shows summary', async ({ browser }) => {
      const adminPage = await browser.newPage();
      await adminPage.goto('/test/login/admin');
      await adminPage.goto('/admin/reports/consultants');

      await expect(adminPage.getByTestId('report-summary')).toBeVisible();

      await adminPage.close();
    });

    test('consultant report shows leaderboard', async ({ browser }) => {
      const adminPage = await browser.newPage();
      await adminPage.goto('/test/login/admin');
      await adminPage.goto('/admin/reports/consultants');

      // Check table headers using role
      await expect(adminPage.getByRole('columnheader', { name: 'Rank' })).toBeVisible();
      await expect(adminPage.getByRole('columnheader', { name: 'Konsultan' })).toBeVisible();
      await expect(adminPage.getByRole('columnheader', { name: 'Enrollments' })).toBeVisible();

      await adminPage.close();
    });

    test('consultant report has filter', async ({ browser }) => {
      const adminPage = await browser.newPage();
      await adminPage.goto('/test/login/admin');
      await adminPage.goto('/admin/reports/consultants');

      await expect(adminPage.getByTestId('report-filter')).toBeVisible();

      await adminPage.close();
    });
  });

  test.describe('Referrer Report', () => {
    test('admin can access referrer report', async ({ browser }) => {
      const adminPage = await browser.newPage();
      await adminPage.goto('/test/login/admin');
      await adminPage.goto('/admin/reports/referrers');

      await expect(adminPage.getByRole('heading', { name: 'Leaderboard Referrer', level: 2 })).toBeVisible();

      await adminPage.close();
    });

    test('referrer report shows summary', async ({ browser }) => {
      const adminPage = await browser.newPage();
      await adminPage.goto('/test/login/admin');
      await adminPage.goto('/admin/reports/referrers');

      // Check summary stats are visible
      await expect(adminPage.getByText('Total Referrer').first()).toBeVisible();
      await expect(adminPage.getByText('Komisi Dibayar').first()).toBeVisible();

      await adminPage.close();
    });

    test('referrer report shows leaderboard', async ({ browser }) => {
      const adminPage = await browser.newPage();
      await adminPage.goto('/test/login/admin');
      await adminPage.goto('/admin/reports/referrers');

      // Check table headers using role
      await expect(adminPage.getByRole('columnheader', { name: 'Referrer' })).toBeVisible();
      await expect(adminPage.getByRole('columnheader', { name: 'Total Referral' })).toBeVisible();
      await expect(adminPage.getByRole('columnheader', { name: 'Komisi Dibayar' })).toBeVisible();

      await adminPage.close();
    });
  });
});
