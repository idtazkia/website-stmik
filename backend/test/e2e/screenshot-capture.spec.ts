import { test, expect, Page } from '@playwright/test';
import * as fs from 'fs';
import * as path from 'path';

/**
 * Screenshot capture for user manual.
 *
 * Each PageDefinition describes a page to screenshot.
 * Screenshots are saved to docs/user-manual/screenshots/.
 *
 * Run: npx playwright test --config=playwright.testrunner.config.ts screenshot-capture
 */

interface PageDefinition {
  id: string;         // screenshot filename (without extension)
  name: string;       // display name in manual
  url: string;        // application route
  section: string;    // subdirectory under screenshots/
  description: string;
  requiresAuth: 'admin' | 'consultant' | 'supervisor' | 'finance' | 'candidate' | 'none';
}

const pages: PageDefinition[] = [
  // --- Admin ---
  { id: 'login', name: 'Login Admin', url: '/admin/login', section: 'admin', description: 'Halaman login admin', requiresAuth: 'none' },
  { id: 'dashboard', name: 'Dashboard Admin', url: '/admin', section: 'admin', description: 'Dashboard admin', requiresAuth: 'admin' },
  { id: 'my-dashboard', name: 'Dashboard EC', url: '/admin/my-dashboard', section: 'admin', description: 'Dashboard Education Consultant', requiresAuth: 'consultant' },
  { id: 'supervisor-dashboard', name: 'Dashboard Supervisor', url: '/admin/supervisor-dashboard', section: 'admin', description: 'Dashboard supervisor', requiresAuth: 'supervisor' },
  { id: 'candidates-list', name: 'Daftar Kandidat', url: '/admin/candidates', section: 'admin', description: 'Halaman daftar kandidat', requiresAuth: 'admin' },
  { id: 'documents-review', name: 'Review Dokumen', url: '/admin/documents', section: 'admin', description: 'Halaman review dokumen', requiresAuth: 'admin' },
  { id: 'campaigns-list', name: 'Daftar Kampanye', url: '/admin/campaigns', section: 'admin', description: 'Halaman daftar kampanye', requiresAuth: 'admin' },
  { id: 'referrers-list', name: 'Daftar Referrer', url: '/admin/referrers', section: 'admin', description: 'Halaman daftar referrer', requiresAuth: 'admin' },
  { id: 'referral-claims', name: 'Referral Claims', url: '/admin/referral-claims', section: 'admin', description: 'Halaman referral claims', requiresAuth: 'admin' },
  { id: 'announcements-list', name: 'Daftar Pengumuman', url: '/admin/announcements', section: 'admin', description: 'Halaman daftar pengumuman', requiresAuth: 'admin' },
  { id: 'report-funnel', name: 'Laporan Funnel', url: '/admin/reports/funnel', section: 'admin', description: 'Laporan analisis funnel', requiresAuth: 'admin' },
  { id: 'report-consultants', name: 'Laporan EC', url: '/admin/reports/consultants', section: 'admin', description: 'Laporan performa Education Consultant', requiresAuth: 'admin' },
  { id: 'report-campaigns', name: 'Laporan Kampanye', url: '/admin/reports/campaigns', section: 'admin', description: 'Laporan efektivitas kampanye', requiresAuth: 'admin' },
  { id: 'report-referrers', name: 'Leaderboard Referrer', url: '/admin/reports/referrers', section: 'admin', description: 'Leaderboard referrer', requiresAuth: 'admin' },
  { id: 'finance-billings', name: 'Daftar Tagihan', url: '/admin/finance/billings', section: 'admin', description: 'Halaman daftar tagihan', requiresAuth: 'finance' },
  { id: 'finance-payments', name: 'Verifikasi Pembayaran', url: '/admin/finance/billings', section: 'admin', description: 'Halaman verifikasi pembayaran', requiresAuth: 'finance' },
  { id: 'settings-users', name: 'Pengaturan User', url: '/admin/settings/users', section: 'admin', description: 'Halaman manajemen user', requiresAuth: 'admin' },
  { id: 'settings-programs', name: 'Program Studi', url: '/admin/settings/programs', section: 'admin', description: 'Halaman pengaturan program studi', requiresAuth: 'admin' },
  { id: 'settings-fees', name: 'Struktur Biaya', url: '/admin/settings/fees', section: 'admin', description: 'Halaman pengaturan biaya', requiresAuth: 'admin' },
  { id: 'settings-assignment', name: 'Algoritma Assignment', url: '/admin/settings/assignment', section: 'admin', description: 'Halaman pengaturan assignment', requiresAuth: 'admin' },
  { id: 'settings-rewards', name: 'Reward Config', url: '/admin/settings/rewards', section: 'admin', description: 'Halaman pengaturan reward', requiresAuth: 'admin' },
  { id: 'settings-document-types', name: 'Jenis Dokumen', url: '/admin/settings/document-types', section: 'admin', description: 'Halaman pengaturan jenis dokumen', requiresAuth: 'admin' },

  // --- Portal ---
  { id: 'login', name: 'Login Portal', url: '/login', section: 'portal', description: 'Halaman login calon mahasiswa', requiresAuth: 'none' },
  { id: 'register-step1', name: 'Registrasi Step 1', url: '/register', section: 'portal', description: 'Form registrasi step 1', requiresAuth: 'none' },
  { id: 'dashboard', name: 'Dashboard Portal', url: '/portal', section: 'portal', description: 'Dashboard calon mahasiswa', requiresAuth: 'candidate' },
  { id: 'documents', name: 'Dokumen Portal', url: '/portal/documents', section: 'portal', description: 'Halaman upload dokumen', requiresAuth: 'candidate' },
  { id: 'payments', name: 'Pembayaran Portal', url: '/portal/payments', section: 'portal', description: 'Halaman status pembayaran', requiresAuth: 'candidate' },
  { id: 'announcements', name: 'Pengumuman Portal', url: '/portal/announcements', section: 'portal', description: 'Halaman pengumuman', requiresAuth: 'candidate' },
  { id: 'referral', name: 'Referral Portal', url: '/portal/referral', section: 'portal', description: 'Halaman program referral', requiresAuth: 'candidate' },
  { id: 'verify-email', name: 'Verifikasi Email', url: '/portal/verify-email', section: 'portal', description: 'Halaman verifikasi email', requiresAuth: 'candidate' },
];

const BASE_URL = process.env.BASE_URL || 'http://localhost:8080';
const SCREENSHOT_DIR = path.join(__dirname, '../../docs/user-manual/screenshots');

async function loginAs(page: Page, role: string) {
  // Uses test login endpoints (only available in test mode)
  if (role === 'candidate') {
    await page.goto(`${BASE_URL}/test/login/candidate`);
  } else {
    await page.goto(`${BASE_URL}/test/login/${role}`);
  }
  await page.waitForLoadState('networkidle');
}

test.describe('User Manual Screenshots', () => {
  test.beforeAll(async () => {
    // Ensure screenshot directories exist
    for (const section of ['admin', 'portal']) {
      const dir = path.join(SCREENSHOT_DIR, section);
      fs.mkdirSync(dir, { recursive: true });
    }
  });

  for (const pageDef of pages) {
    test(`capture: ${pageDef.section}/${pageDef.id}`, async ({ page }) => {
      page.setViewportSize({ width: 1280, height: 800 });

      if (pageDef.requiresAuth !== 'none') {
        await loginAs(page, pageDef.requiresAuth);
      }

      await page.goto(`${BASE_URL}${pageDef.url}`);
      await page.waitForLoadState('networkidle');

      const screenshotPath = path.join(SCREENSHOT_DIR, pageDef.section, `${pageDef.id}.png`);
      await page.screenshot({ path: screenshotPath, fullPage: false });

      expect(fs.existsSync(screenshotPath)).toBeTruthy();
    });
  }

  // Candidate detail: navigate from list to get a real candidate ID
  test('capture: admin/candidate-detail', async ({ page }) => {
    page.setViewportSize({ width: 1280, height: 800 });
    await loginAs(page, 'admin');

    await page.goto(`${BASE_URL}/admin/candidates`);
    await page.waitForLoadState('networkidle');

    const firstLink = page.locator('[data-testid^="view-candidate-"]').first();
    await firstLink.click();
    await page.waitForLoadState('networkidle');

    const screenshotPath = path.join(SCREENSHOT_DIR, 'admin', 'candidate-detail.png');
    await page.screenshot({ path: screenshotPath, fullPage: false });
    expect(fs.existsSync(screenshotPath)).toBeTruthy();
  });
});
