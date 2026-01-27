import { test, expect } from '@playwright/test';

test.describe('Commission Management', () => {
  test.describe('Commission List', () => {
    test('admin can access commissions page', async ({ browser }) => {
      const adminPage = await browser.newPage();
      await adminPage.goto('/test/login/admin');
      await adminPage.goto('/admin/commissions');

      await expect(adminPage.getByTestId('commissions-page')).toBeVisible();

      await adminPage.close();
    });

    test('finance user can access commissions page', async ({ browser }) => {
      const financePage = await browser.newPage();
      await financePage.goto('/test/login/finance');
      await financePage.goto('/admin/commissions');

      await expect(financePage.getByTestId('commissions-page')).toBeVisible();

      await financePage.close();
    });

    test('commissions page shows stats', async ({ browser }) => {
      const adminPage = await browser.newPage();
      await adminPage.goto('/test/login/admin');
      await adminPage.goto('/admin/commissions');

      // Check stats are visible
      await expect(adminPage.locator('.grid.grid-cols-3')).toBeVisible();

      await adminPage.close();
    });

    test('commissions page has filter tabs', async ({ browser }) => {
      const adminPage = await browser.newPage();
      await adminPage.goto('/test/login/admin');
      await adminPage.goto('/admin/commissions');

      await expect(adminPage.getByTestId('filter-tabs')).toBeVisible();
      await expect(adminPage.getByTestId('filter-pending')).toBeVisible();
      await expect(adminPage.getByTestId('filter-approved')).toBeVisible();
      await expect(adminPage.getByTestId('filter-paid')).toBeVisible();

      await adminPage.close();
    });
  });

  test.describe('Commission Filters', () => {
    test('filter by pending status', async ({ browser }) => {
      const adminPage = await browser.newPage();
      await adminPage.goto('/test/login/admin');
      await adminPage.goto('/admin/commissions');

      await adminPage.getByTestId('filter-pending').click();
      await expect(adminPage).toHaveURL(/status=pending/);

      await adminPage.close();
    });

    test('filter by approved status', async ({ browser }) => {
      const adminPage = await browser.newPage();
      await adminPage.goto('/test/login/admin');
      await adminPage.goto('/admin/commissions');

      await adminPage.getByTestId('filter-approved').click();
      await expect(adminPage).toHaveURL(/status=approved/);

      await adminPage.close();
    });

    test('filter by paid status', async ({ browser }) => {
      const adminPage = await browser.newPage();
      await adminPage.goto('/test/login/admin');
      await adminPage.goto('/admin/commissions');

      await adminPage.getByTestId('filter-paid').click();
      await expect(adminPage).toHaveURL(/status=paid/);

      await adminPage.close();
    });
  });

  test.describe('Commission CSV Export', () => {
    test('export button is visible on commissions page', async ({ browser }) => {
      const adminPage = await browser.newPage();
      await adminPage.goto('/test/login/admin');
      await adminPage.goto('/admin/commissions');

      await expect(adminPage.getByTestId('btn-export')).toBeVisible();
      await expect(adminPage.getByTestId('btn-export')).toContainText('Export untuk Transfer');

      await adminPage.close();
    });

    test('export button links to correct endpoint', async ({ browser }) => {
      const adminPage = await browser.newPage();
      await adminPage.goto('/test/login/admin');
      await adminPage.goto('/admin/commissions');

      const exportLink = adminPage.getByTestId('btn-export');
      const href = await exportLink.getAttribute('href');
      expect(href).toBe('/admin/commissions/export?status=approved');

      await adminPage.close();
    });

    test('CSV export returns proper response headers', async ({ browser }) => {
      const adminPage = await browser.newPage();
      await adminPage.goto('/test/login/admin');

      // Use page.request to fetch file directly (avoids download dialog)
      const response = await adminPage.request.get('/admin/commissions/export?status=approved');

      expect(response.status()).toBe(200);
      expect(response.headers()['content-type']).toContain('text/csv');
      expect(response.headers()['content-disposition']).toContain('attachment');
      expect(response.headers()['content-disposition']).toContain('.csv');

      await adminPage.close();
    });

    test('CSV export contains proper headers', async ({ browser }) => {
      const adminPage = await browser.newPage();
      await adminPage.goto('/test/login/admin');

      // Use page.request to fetch file directly (avoids download dialog)
      const response = await adminPage.request.get('/admin/commissions/export?status=approved');
      const body = await response.text();

      // Check CSV header row (after BOM)
      expect(body).toContain('No,Nama Referrer,Tipe,Nama Bank,No Rekening,Atas Nama,Jumlah,Kandidat,Trigger Event,Tanggal Approve');

      await adminPage.close();
    });

    test('finance user can export commissions', async ({ browser }) => {
      const financePage = await browser.newPage();
      await financePage.goto('/test/login/finance');

      // Use page.request to fetch file directly (avoids download dialog)
      const response = await financePage.request.get('/admin/commissions/export?status=approved');
      expect(response.status()).toBe(200);
      expect(response.headers()['content-type']).toContain('text/csv');

      await financePage.close();
    });

    test('consultant cannot export commissions', async ({ browser }) => {
      const consultantPage = await browser.newPage();
      await consultantPage.goto('/test/login/consultant');

      const response = await consultantPage.goto('/admin/commissions/export?status=approved');
      expect(response?.status()).toBe(403);

      await consultantPage.close();
    });
  });

  test.describe('Commission Actions', () => {
    test('batch approve button exists and is disabled when no selection', async ({ browser }) => {
      const adminPage = await browser.newPage();
      await adminPage.goto('/test/login/admin');
      await adminPage.goto('/admin/commissions');

      const batchApproveBtn = adminPage.getByTestId('btn-batch-approve');
      await expect(batchApproveBtn).toBeVisible();
      await expect(batchApproveBtn).toBeDisabled();

      await adminPage.close();
    });
  });
});
