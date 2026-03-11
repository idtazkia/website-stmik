import { test, expect } from '@playwright/test';

test.describe('Admin Dashboard', () => {
  test('dashboard loads with all sections visible', async ({ browser }) => {
    const page = await browser.newPage();
    await page.goto('/test/login/admin');
    await page.goto('/admin');

    await expect(page.getByTestId('admin-dashboard')).toBeVisible();
    await expect(page.getByTestId('stats-cards')).toBeVisible();
    await expect(page.getByTestId('overdue-section')).toBeVisible();
    await expect(page.getByTestId('today-tasks-section')).toBeVisible();
    await expect(page.getByTestId('funnel-section')).toBeVisible();
    await expect(page.getByTestId('recent-candidates-section')).toBeVisible();

    await page.close();
  });

  test('stats cards show numeric values', async ({ browser }) => {
    const page = await browser.newPage();
    await page.goto('/test/login/admin');
    await page.goto('/admin');

    const statsCards = page.getByTestId('stats-cards');

    // Verify stat labels are present
    await expect(statsCards.getByText('Total Kandidat')).toBeVisible();
    await expect(statsCards.getByText('Prospecting')).toBeVisible();
    await expect(statsCards.getByText('Committed')).toBeVisible();
    await expect(statsCards.getByText('Enrolled')).toBeVisible();

    // Verify each card has a bold numeric value (text-3xl font-bold)
    const boldValues = statsCards.locator('.text-3xl.font-bold');
    const count = await boldValues.count();
    expect(count).toBe(4);

    for (let i = 0; i < count; i++) {
      const text = await boldValues.nth(i).textContent();
      expect(text).toBeTruthy();
      expect(text!.trim()).toMatch(/^\d+$/);
    }

    await page.close();
  });

  test('overdue section renders correctly', async ({ browser }) => {
    const page = await browser.newPage();
    await page.goto('/test/login/admin');
    await page.goto('/admin');

    const overdueSection = page.getByTestId('overdue-section');
    await expect(overdueSection).toBeVisible();

    // Check heading
    await expect(overdueSection.getByRole('heading', { name: 'Follow-up Terlambat' })).toBeVisible();

    // Check badge shows candidate count
    const badge = overdueSection.locator('.bg-red-100');
    await expect(badge).toBeVisible();
    await expect(badge).toContainText('kandidat');

    // If there are overdue candidates, verify items have link to detail page
    const overdueItems = overdueSection.locator('.bg-red-50');
    const itemCount = await overdueItems.count();
    if (itemCount > 0) {
      for (let i = 0; i < itemCount; i++) {
        const item = overdueItems.nth(i);
        // Each item should have a "Lihat" link pointing to candidate detail
        const link = item.getByRole('link', { name: 'Lihat' });
        await expect(link).toBeVisible();
        await expect(link).toHaveAttribute('href', /\/admin\/candidates\//);
      }
    } else {
      // Empty state message
      await expect(overdueSection.getByText('Tidak ada follow-up terlambat')).toBeVisible();
    }

    await page.close();
  });

  test('today tasks section renders correctly', async ({ browser }) => {
    const page = await browser.newPage();
    await page.goto('/test/login/admin');
    await page.goto('/admin');

    const todaySection = page.getByTestId('today-tasks-section');
    await expect(todaySection).toBeVisible();

    // Check heading
    await expect(todaySection.getByRole('heading', { name: 'Tugas Hari Ini' })).toBeVisible();

    // Check badge shows follow-up count
    const badge = todaySection.locator('.bg-blue-100');
    await expect(badge).toBeVisible();
    await expect(badge).toContainText('follow-up');

    // If there are tasks, verify items have link to detail page
    const taskItems = todaySection.locator('.bg-gray-50.rounded-lg');
    const itemCount = await taskItems.count();
    if (itemCount > 0) {
      for (let i = 0; i < itemCount; i++) {
        const item = taskItems.nth(i);
        // Each item should have a "Follow-up" link pointing to candidate detail
        const link = item.getByRole('link', { name: 'Follow-up' });
        await expect(link).toBeVisible();
        await expect(link).toHaveAttribute('href', /\/admin\/candidates\//);
      }
    } else {
      // Empty state message
      await expect(todaySection.getByText('Tidak ada tugas hari ini')).toBeVisible();
    }

    await page.close();
  });

  test('funnel overview shows stages with labels', async ({ browser }) => {
    const page = await browser.newPage();
    await page.goto('/test/login/admin');
    await page.goto('/admin');

    const funnelSection = page.getByTestId('funnel-section');
    await expect(funnelSection).toBeVisible();

    // Check heading
    await expect(funnelSection.getByText('Funnel Overview')).toBeVisible();

    // Check all 4 funnel stage labels
    await expect(funnelSection.getByText('Registered')).toBeVisible();
    await expect(funnelSection.getByText('Prospecting')).toBeVisible();
    await expect(funnelSection.getByText('Committed')).toBeVisible();
    await expect(funnelSection.getByText('Enrolled')).toBeVisible();

    // Check that each stage has a numeric value (text-2xl font-bold)
    const stageValues = funnelSection.locator('.text-2xl.font-bold');
    const count = await stageValues.count();
    expect(count).toBe(4);

    for (let i = 0; i < count; i++) {
      const text = await stageValues.nth(i).textContent();
      expect(text).toBeTruthy();
      expect(text!.trim()).toMatch(/^\d+$/);
    }

    await page.close();
  });

  test('recent candidates table shows headers and data', async ({ browser }) => {
    const page = await browser.newPage();
    await page.goto('/test/login/admin');
    await page.goto('/admin');

    const recentSection = page.getByTestId('recent-candidates-section');
    await expect(recentSection).toBeVisible();

    // Check section heading
    await expect(recentSection.getByText('Kandidat Terbaru')).toBeVisible();

    // Check "Lihat Semua" link
    const viewAllLink = recentSection.getByRole('link', { name: /Lihat Semua/ });
    await expect(viewAllLink).toBeVisible();
    await expect(viewAllLink).toHaveAttribute('href', '/admin/candidates');

    // Check table headers
    await expect(recentSection.getByRole('columnheader', { name: /nama/i })).toBeVisible();
    await expect(recentSection.getByRole('columnheader', { name: /prodi/i })).toBeVisible();
    await expect(recentSection.getByRole('columnheader', { name: /status/i })).toBeVisible();
    await expect(recentSection.getByRole('columnheader', { name: /konsultan/i })).toBeVisible();
    await expect(recentSection.getByRole('columnheader', { name: /tanggal/i })).toBeVisible();

    // If there are candidate rows, verify each has a link to detail page
    const dataRows = recentSection.locator('tbody tr');
    const rowCount = await dataRows.count();
    if (rowCount > 0) {
      const firstRowLink = dataRows.first().getByRole('link');
      const firstRowLinkCount = await firstRowLink.count();
      if (firstRowLinkCount > 0) {
        await expect(firstRowLink.first()).toHaveAttribute('href', /\/admin\/candidates\//);
      }
    }

    await page.close();
  });

  test('consultant accessing /admin sees appropriate dashboard', async ({ browser }) => {
    const page = await browser.newPage();
    await page.goto('/test/login/consultant');
    await page.goto('/admin');

    // Consultant should see the admin dashboard (their own view)
    await expect(page.getByTestId('admin-dashboard')).toBeVisible();
    await expect(page.getByTestId('stats-cards')).toBeVisible();

    await page.close();
  });

  test('supervisor accessing /admin sees admin dashboard', async ({ browser }) => {
    const page = await browser.newPage();
    await page.goto('/test/login/supervisor');
    await page.goto('/admin');

    // Supervisor should see the admin dashboard
    await expect(page.getByTestId('admin-dashboard')).toBeVisible();
    await expect(page.getByTestId('stats-cards')).toBeVisible();
    await expect(page.getByTestId('funnel-section')).toBeVisible();
    await expect(page.getByTestId('recent-candidates-section')).toBeVisible();

    await page.close();
  });
});
